package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/soerenschneider/hermes/internal"
	"github.com/soerenschneider/hermes/internal/config"
	"github.com/soerenschneider/hermes/internal/events"
	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/internal/notification"
	"github.com/soerenschneider/hermes/internal/queue"
	"github.com/soerenschneider/hermes/internal/validation"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	defaultConfigFile = "/etc/hermes.yaml"
)

var (
	configFile   string
	printVersion bool
	debug        bool
)

type Deps struct {
	conf config.Config

	queue  queue.Queue
	cortex *notification.Dispatcher

	wg *sync.WaitGroup

	eventSources []events.EventSource
}

func main() {
	log.Info().Msgf("Starting hermes, version %s (%s)", internal.BuildVersion, internal.CommitHash)

	// parse flags
	parseFlags()
	if printVersion { // abusing bool as subcmd
		fmt.Println(internal.BuildVersion)
		os.Exit(0)
	}

	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// parse config
	conf, err := config.Read(configFile)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read config")
	}
	log.Info().Msg("Read config")

	if err := validation.Validate(conf); err != nil {
		log.Fatal().Err(err).Msg("invalid config")
	}
	log.Info().Msg("Validation passed successfully")

	// build deps
	deps := &Deps{}
	deps.conf = *conf
	deps.wg = &sync.WaitGroup{}

	deps.queue, err = deps.buildQueue(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build queue")
	}

	deps.cortex, err = deps.buildCortex(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build cortex")
	}

	log.Info().Msg("Building event sources")
	deps.eventSources, err = buildEventSources(conf)
	if err != nil {
		log.Fatal().Err(err).Msg("could not build event sources")
	}

	run(deps)
}

func run(deps *Deps) {
	ctx, cancel := context.WithCancel(context.Background())

	for i := 0; i < 10; i++ {
		go deps.cortex.Listen(ctx)
	}
	go deps.cortex.StartQueueReconciliation(ctx)

	for _, eventSource := range deps.eventSources {
		go func(source events.EventSource) {
			err := source.Listen(ctx, deps.cortex, deps.wg)
			if err != nil {
				log.Error().Err(err).Msg("listening on event source failed")
			}
		}(eventSource)
	}

	if len(deps.conf.MetricsAddr) > 0 {
		go metrics.StartServer(deps.conf.MetricsAddr)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	<-sig
	log.Info().Msg("Caught signal, shutting down gracefully")
	cancel()

	deps.wg.Wait()
	log.Info().Msg("All components shut down, bye!")
}

func parseFlags() {
	flag.StringVar(&configFile, "config", defaultConfigFile, fmt.Sprintf("Path to the config file (default %s)", defaultConfigFile))
	flag.BoolVar(&printVersion, "version", false, "Print printVersion and exit")
	flag.BoolVar(&debug, "debug", false, "Print debug information")
	flag.Parse()
}
