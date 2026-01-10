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
	telemetryAddr     = flag.String("addr", ":9101", "Address for pwrstat exporter")
	pwrstatStatusPath = flag.String("pwrstat-status-path", "/var/lib/pwrstat_status/status", "Path to pwrstat -status output")
)

func main() {
	flag.Parse()

	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
	)

	pwrstatReader := pwrstat.NewReader(*pwrstatStatusPath)
	isPwrstatInstalled := pwrstatReader.IsExist()
	if !isPwrstatInstalled {
		logger.Error("Please check if the pwrstat -status output file exists", slog.String("path", *pwrstatStatusPath))
		os.Exit(1)
	}

	collector := collector.New(logger, pwrstatReader)
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
