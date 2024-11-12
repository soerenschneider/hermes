package sqlite

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/soerenschneider/hermes/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestNewSQLiteQueue(t *testing.T) {
	file, err := os.CreateTemp("/tmp", "hermes-sqlite")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(file.Name())
	}()

	db, err := New(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	item := domain.Notification{
		ServiceId: "svc",
		Subject:   "subject",
		Message:   "message message message",
	}
	ctx := context.Background()

	isEmpty, err := db.IsEmpty(ctx)
	assert.NoError(t, err)
	assert.Equal(t, true, isEmpty)

	if err := db.Offer(ctx, item); err != nil {
		t.Fatal(err)
	}

	isEmpty, err = db.IsEmpty(ctx)
	assert.NoError(t, err)
	assert.Equal(t, false, isEmpty)

	read, err := db.Get(ctx)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, item.Message, read.Message)
	assert.Equal(t, item.Subject, read.Subject)
	assert.Equal(t, item.ServiceId, read.ServiceId)
	assert.True(t, time.Since(read.Inserted) < 5*time.Second)

	isEmpty, err = db.IsEmpty(ctx)
	assert.NoError(t, err)
	assert.Equal(t, true, isEmpty)
}
