package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/internal/queue/sqlite/generated"
	"github.com/soerenschneider/hermes/pkg"
)

type SQLiteQueue struct {
	db        *sql.DB
	generated *generated.Queries
}

func New(dbPath string) (*SQLiteQueue, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	gen := generated.New(db)
	ret := &SQLiteQueue{
		db:        db,
		generated: gen,
	}

	return ret, ret.Migrate(context.Background())
}

func (q *SQLiteQueue) Offer(ctx context.Context, item pkg.Notification) error {
	params := generated.InsertParams{
		Subject:   item.Subject,
		Message:   item.Message,
		ServiceID: item.ServiceId,
		Retries:   int64(item.UnsuccessfulAttempts),
		RetryDate: item.RetryDate,
	}
	return q.generated.Insert(ctx, params)
}

func (q *SQLiteQueue) Get(ctx context.Context) (pkg.Notification, error) {
	read, err := q.generated.GetMessage(ctx)
	if err != nil {
		return pkg.Notification{}, err
	}

	return pkg.Notification{
		Id:                   read.ID,
		Inserted:             read.InsertionDate,
		UnsuccessfulAttempts: int(read.Retries),
		ServiceId:            read.ServiceID,
		Subject:              read.Subject,
		Message:              read.Message,
	}, nil
}

func (q *SQLiteQueue) GetMessageCount(ctx context.Context) (int64, error) {
	cnt, err := q.generated.GetCount(ctx)
	if err == nil {
		metrics.QueueSize.Set(float64(cnt))
	}
	return cnt, err
}

func (q *SQLiteQueue) IsEmpty(ctx context.Context) (bool, error) {
	cnt, err := q.generated.GetCount(ctx)
	return cnt == 0, err
}

func (db *SQLiteQueue) Migrate(ctx context.Context) error {
	if schemaVersionReadError != nil {
		return schemaVersionReadError
	}

	var currentVersion int
	_ = db.db.QueryRowContext(ctx, `SELECT version FROM schema_version`).Scan(&currentVersion)

	log.Info().Msgf("Current DB schema at version %d, latest schema version is %d", currentVersion, schemaVersion)
	if currentVersion >= schemaVersion {
		return nil
	}

	migrations, err := GetMigrations()
	if err != nil {
		return err
	}

	for version := currentVersion; version < schemaVersion; version++ {
		newVersion := version + 1

		tx, err := db.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("can not start transaction %w", err)
		}

		sql := migrations[version]
		_, err = tx.ExecContext(ctx, string(sql))
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}

		if _, err := tx.ExecContext(ctx, `DELETE FROM schema_version`); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_version (version) VALUES ($1)`, newVersion); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("[Migration v%d] %v", newVersion, err)
		}
		log.Info().Msgf("Successfully migrated DB to version %d", newVersion)
	}

	return nil
}
