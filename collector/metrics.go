package collector

import "github.com/prometheus/client_golang/prometheus"

var (
	// Info metrics
	upsProperties = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_properties"),
		"UPS model and configuration information",
		[]string{"model", "firmware", "rating_voltage", "rating_power"},
		nil,
	)

	// Gauge metrics
	utilityVoltage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_utility_voltage_volts"),
		"Input voltage from utility power",
		nil,
		nil,
	)

	outputVoltage = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_output_voltage_volts"),
		"Output voltage from UPS",
		nil,
		nil,
	)

	batteryCapacity = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_battery_cap_percentages"),
		"UPS battery capacity",
		nil,
		nil,
	)

	batteryRuntime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_battery_runtime_minutes"),
		"UPS battery remaining runtime",
		nil,
		nil,
	)

	loadWatts = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_load_watts"),
		"UPS load in watts",
		nil,
		nil,
	)

	loadPercent = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_load_percentages"),
		"UPS load in percentages",
		nil,
		nil,
	)

	// Enum metrics
	state = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_state"),
		"Current UPS state",
		[]string{"state"},
		nil,
	)

	powerSource = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_power_source"),
		"Current power source",
		[]string{"source"},
		nil,
	)

	// Counter metrics
	powerEventsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_power_events_total"),
		"Total number of power events",
		[]string{"type"},
		nil,
	)

	batteryTestsTotal = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_battery_tests_total"),
		"Total number of battery tests",
		[]string{"result"},
		nil,
	)

	// Timestamp metrics
	lastPowerEventTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_last_power_event_timestamp_seconds"),
		"Timestamp of the last power event",
		[]string{"type"},
		nil,
	)

	lastBatteryTestTime = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "ups_last_battery_test_timestamp_seconds"),
		"Timestamp of the last battery test",
		[]string{"result"},
		nil,
	)
)
