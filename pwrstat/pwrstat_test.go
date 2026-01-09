package pwrstat

import (
	"fmt"
)

func Example_parse() {
	rawStatus := []byte(`

	The UPS information shows as following:

		Properties:
			Model Name................... CP900EPFCLCD
			Firmware Number.............. BF02405AAG1
			Rating Voltage............... 230 V
			Rating Power................. 540 Watt

		Current UPS status:
			State........................ Normal
			Power Supply by.............. Utility Power
			Utility Voltage.............. 232 V
			Output Voltage............... 232 V
			Battery Capacity............. 100 %
			Remaining Runtime............ 48 min.
			Load......................... 54 Watt(10 %)
			Line Interaction............. None
			Test Result.................. Passed at 2026/01/07 13:54:02
			Last Power Event............. Blackout at 2026/01/09 00:50:35 for 1 min.

		`)

	status := parse(string(rawStatus))

	fmt.Printf(
		"Model Name %s\n"+
			"Firmware Number %s\n"+
			"Rating Voltage %s\n"+
			"Rating Power %s\n"+

			"State %s\n"+
			"Power Supply by %s\n"+
			"Utility Voltage %s\n"+
			"Output Voltage %s\n"+
			"Battery Capacity %s\n"+
			"Remaining Runtime %s\n"+
			"Load %s\n"+
			"Line Interaction %s\n"+
			"Test Result %s\n"+
			"Last Power Event %s\n",

		status.ModelName,
		status.FirmwareNumber,
		status.RatingVoltage,
		status.RatingPower,

		status.State,
		status.PowerSupplyBy,
		status.UtilityVoltage,
		status.OutputVoltage,
		status.BatteryCapacity,
		status.RemainingRuntime,
		status.Load,
		status.LineInteraction,
		status.TestResult,
		status.LastPowerEvent,
	)

	// Output:
	// Model Name CP900EPFCLCD
	// Firmware Number BF02405AAG1
	// Rating Voltage 230 V
	// Rating Power 540 Watt
	// State Normal
	// Power Supply by Utility Power
	// Utility Voltage 232 V
	// Output Voltage 232 V
	// Battery Capacity 100 %
	// Remaining Runtime 48 min.
	// Load 54 Watt(10 %)
	// Line Interaction None
	// Test Result Passed at 2026/01/07 13:54:02
	// Last Power Event Blackout at 2026/01/09 00:50:35 for 1 min.
}
