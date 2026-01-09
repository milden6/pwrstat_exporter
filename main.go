package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/milden6/pwrstat_exporter/collector"
	"github.com/milden6/pwrstat_exporter/pwrstat"
	"github.com/milden6/pwrstat_exporter/server"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9101", "Address for pwrstat exporter")
)

func main() {
	flag.Parse()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	isPwrstatInstalled := pwrstat.IsExist()
	if !isPwrstatInstalled {
		logger.Error("Please install pwrstat first")
		os.Exit(1)
	}

	collector := collector.New(logger)
	prometheus.MustRegister(collector)

	server := server.New(*telemetryAddr)

	go func() {
		logger.Info("Server starting", slog.String("port", *telemetryAddr))
		if err := server.Start(); err != nil {
			panic(err)
		}
	}()

	rootCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-rootCtx.Done()
	logger.Info("Server stopped gracefully")
}
