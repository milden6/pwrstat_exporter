package pwrstat

import (
	"os/exec"
	"strings"
)

const (
	modelName        = "Model Name"
	firmwareNumber   = "Firmware Number"
	ratingVoltage    = "Rating Voltage"
	ratingPower      = "Rating Power"
	state            = "State"
	powerSupplyBy    = "Power Supply by"
	utilityVoltage   = "Utility Voltage"
	outputVoltage    = "Output Voltage"
	batteryCapacity  = "Battery Capacity"
	remainingRuntime = "Remaining Runtime"
	load             = "Load"
	lineInteraction  = "Line Interaction"
	testResult       = "Test Result"
	lastPowerEvent   = "Last Power Event"
)

type UPSStatus struct {
	// UPS properties
	ModelName      string
	FirmwareNumber string
	RatingVoltage  string
	RatingPower    string

	// Current UPS status
	State            string
	PowerSupplyBy    string
	UtilityVoltage   string
	OutputVoltage    string
	BatteryCapacity  string
	RemainingRuntime string
	Load             string
	LineInteraction  string
	TestResult       string
	LastPowerEvent   string
}

func IsExist() bool {
	_, err := exec.LookPath("pwrstat")

	return err == nil
}

func Status() (*UPSStatus, error) {
	rawStatus, err := exec.Command("pwrstat", "-status").Output()
	if err != nil {
		return nil, err
	}

	return parse(string(rawStatus)), nil
}

func parse(rawStatus string) *UPSStatus {
	upsStatus := &UPSStatus{}

	lines := strings.SplitSeq(string(rawStatus), "\n")
	strReplacer := strings.NewReplacer("\t", "", ".", "")

	for line := range lines {
		keyValue := strings.Split(line, ". ")
		if len(keyValue) == 2 {
			key := strReplacer.Replace(keyValue[0])
			value := keyValue[1]

			switch key {
			case modelName:
				upsStatus.ModelName = value
			case firmwareNumber:
				upsStatus.FirmwareNumber = value
			case ratingVoltage:
				upsStatus.RatingVoltage = value
			case ratingPower:
				upsStatus.RatingPower = value
			case state:
				upsStatus.State = value
			case powerSupplyBy:
				upsStatus.PowerSupplyBy = value
			case utilityVoltage:
				upsStatus.UtilityVoltage = value
			case outputVoltage:
				upsStatus.OutputVoltage = value
			case batteryCapacity:
				upsStatus.BatteryCapacity = value
			case remainingRuntime:
				upsStatus.RemainingRuntime = value
			case load:
				upsStatus.Load = value
			case lineInteraction:
				upsStatus.LineInteraction = value
			case testResult:
				upsStatus.TestResult = value
			case lastPowerEvent:
				upsStatus.LastPowerEvent = value

			}
		}
	}

	return upsStatus
}
