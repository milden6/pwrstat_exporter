package collector

import (
	"log/slog"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/milden6/pwrstat_exporter/pwrstat"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "pwrstat"
)

var reTimestamp = regexp.MustCompile(`\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2}`)

type Collector struct {
	logger        *slog.Logger
	pwrstatReader *pwrstat.Reader
}

func New(logger *slog.Logger, pr *pwrstat.Reader) *Collector {
	return &Collector{
		logger:        logger,
		pwrstatReader: pr,
	}
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- upsProperties
	ch <- utilityVoltage
	ch <- outputVoltage
	ch <- batteryCapacity
	ch <- batteryRuntime
	ch <- loadWatts
	ch <- loadPercent
	ch <- state
	ch <- powerSource
	ch <- powerEventsTotal
	ch <- batteryTestsTotal
	ch <- lastPowerEventTime
	ch <- lastBatteryTestTime
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	status, err := c.pwrstatReader.Status()
	if err != nil {
		c.logger.Error("failed to get UPS status", slog.Any("error", err))

		return
	}

	ch <- prometheus.MustNewConstMetric(
		upsProperties,
		prometheus.GaugeValue,
		1,
		status.ModelName,
		status.FirmwareNumber,
		status.RatingVoltage,
		status.RatingPower,
	)

	ch <- prometheus.MustNewConstMetric(
		utilityVoltage,
		prometheus.GaugeValue,
		c.strToFloat(status.UtilityVoltage),
	)

	ch <- prometheus.MustNewConstMetric(
		outputVoltage,
		prometheus.GaugeValue,
		c.strToFloat(status.OutputVoltage),
	)

	ch <- prometheus.MustNewConstMetric(
		batteryCapacity,
		prometheus.GaugeValue,
		c.strToFloat(status.BatteryCapacity),
	)

	ch <- prometheus.MustNewConstMetric(
		batteryRuntime,
		prometheus.GaugeValue,
		c.strToFloat(status.RemainingRuntime),
	)

	numLoadWatts := c.strToFloat(status.Load)
	ch <- prometheus.MustNewConstMetric(
		loadWatts,
		prometheus.GaugeValue,
		numLoadWatts,
	)

	numLoadPercenatges := (numLoadWatts / c.strToFloat(status.RatingPower)) * 100
	ch <- prometheus.MustNewConstMetric(
		loadPercent,
		prometheus.GaugeValue,
		numLoadPercenatges,
	)

	states := []string{"normal", "power failure"}
	for _, s := range states {
		value := 0.0

		if strings.ToLower(status.State) == s {
			value = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			state,
			prometheus.GaugeValue,
			value,
			s,
		)
	}

	sources := []string{"utility power", "battery power"}
	for _, src := range sources {
		value := 0.0

		if strings.ToLower(status.PowerSupplyBy) == src {
			value = 1.0
		}

		ch <- prometheus.MustNewConstMetric(
			powerSource,
			prometheus.GaugeValue,
			value,
			src,
		)
	}

	// TODO: add counter metrics

	// ch <- prometheus.MustNewConstMetric(
	// 	powerEventsTotal,
	// 	prometheus.CounterValue,
	// 	float64(blackoutCount),
	// 	"blackout",
	// )

	// ch <- prometheus.MustNewConstMetric(
	// 	batteryTestsTotal,
	// 	prometheus.CounterValue,
	// 	float64(testsPassed),
	// 	"passed",
	// )

	eventTimeStr := reTimestamp.FindString(status.LastPowerEvent)
	eventTime, err := time.Parse("2006/01/02 15:04:05", eventTimeStr)
	if err != nil {
		c.logger.Error("failed to parse time", slog.Any("error", err))

		return
	}

	ch <- prometheus.MustNewConstMetric(
		lastPowerEventTime,
		prometheus.GaugeValue,
		float64(eventTime.Unix()),
		"blackout",
	)

	testTimeStr := reTimestamp.FindString(status.TestResult)
	testTime, err := time.Parse("2006/01/02 15:04:05", testTimeStr)
	if err != nil {
		c.logger.Error("failed to parse time", slog.Any("error", err))

		return
	}

	ch <- prometheus.MustNewConstMetric(
		lastBatteryTestTime,
		prometheus.GaugeValue,
		float64(testTime.Unix()),
		"passed",
	)
}

func (c *Collector) strToFloat(s string) float64 {
	elems := strings.Split(s, " ")
	num, err := strconv.Atoi(elems[0])
	if err != nil {
		c.logger.Error("failed Atoi", slog.Any("error", err))

		return 0
	}

	return float64(num)
}
