package pwacii

import (
	"fmt"
)

type EntranceStatus string

// X Constants - Entrance Status
const XN EntranceStatus = "HalfBodyTurnstileEntryNotReleased"
const XS EntranceStatus = "HalfBodyTurnstileEntryReleased"
const XI EntranceStatus = "UndefinedEntryStatus"

type ExitStatus string

// Y Constants - Exit Status
const YN ExitStatus = "HalfBodyTurnstileExitNotReleased"
const YS ExitStatus = "HalfBodyTurnstileExitReleased"
const YI ExitStatus = "UndefinedExitStatus"

type MovementStatus string

// Z Constants - Movement Status
const ZP MovementStatus = "HalfBodyTurnstileAtRest"
const ZE MovementStatus = "HalfBodyTurnstileEntryMotion"
const ZS MovementStatus = "HalfBodyTurnstileExitMotion"
const ZI MovementStatus = "UndefinedMotionStatus"

type MovementProgress string

// W Constants - Rest Status
const W0 MovementProgress = "HalfBodyTurnstileAtRest"
const W1 MovementProgress = "HalfBodyTurnstileFirstQuarterMotion"
const W2 MovementProgress = "HalfBodyTurnstileMidRotation"
const W3 MovementProgress = "HalfBodyTurnstileThirdQuarterMotion"
const WI MovementProgress = "UndefinedWStatus"

type Status struct {
	EntranceStatus   EntranceStatus   // X
	ExitStatus       ExitStatus       // Y
	MovementStatus   MovementStatus   // Z
	MovementProgress MovementProgress // W
}

func ParseStatus(data string) (Status, error) {
	if len(data) < 4 {
		return Status{}, fmt.Errorf("invalid data length, %q", data)
	}

	var status Status
	switch data[0] {
	case 'N':
		status.EntranceStatus = XN
	case 'S':
		status.EntranceStatus = XS
	case 'I':
		status.EntranceStatus = XI
	default:
		return Status{}, fmt.Errorf("invalid entrance status, %q", data)
	}

	switch data[1] {
	case 'N':
		status.ExitStatus = YN
	case 'S':
		status.ExitStatus = YS
	case 'I':
		status.ExitStatus = YI
	default:
		return Status{}, fmt.Errorf("invalid exit status, %q", data)
	}

	switch data[2] {
	case 'P':
		status.MovementStatus = ZP
	case 'E':
		status.MovementStatus = ZE
	case 'S':
		status.MovementStatus = ZS
	case 'I':
		status.MovementStatus = ZI
	default:
		return Status{}, fmt.Errorf("invalid movement status, %q", data)
	}

	switch data[3] {
	case '0':
		status.MovementProgress = W0
	case '1':
		status.MovementProgress = W1
	case '2':
		status.MovementProgress = W2
	case '3':
		status.MovementProgress = W3
	case 'I':
		status.MovementProgress = WI
	default:
		return Status{}, fmt.Errorf("invalid movement progress, %q", data)
	}

	return status, nil
}
