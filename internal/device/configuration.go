package device

import (
	"fmt"
	"strings"
)

// Inicialmente, establezco una estructura llamada Opts que tendrá todos los tipos como campos.
type Opts struct {
	Lock                    Lock
	RegisterType            RegisterType
	ControlType             ControlType
	PictoMapType            PictoMapType
	ViolationControlType    ViolationControlType
	TurnstileType           TurnstileType
	ReleaseControlType      ReleaseControlType
	BQCType                 BQCType
	CheckType               CheckType
	CheckPercent            CheckPercent
	TimeoutEntry            TimeoutEntry
	MaxAccumEntry           MaxAccumEntry
	TimeoutSelenoideControl TimeoutSelenoideControlCorredera
	OneButtonsControl       ButtonsControl
	TwoButtonsControl       ButtonsControl
	CorrederaType           CorrederaType
	PictoControlType        PictoControlType
}

func (opt *Opts) OptsToString() string {
	s := strings.Builder{}

	// Lock
	switch opt.Lock {
	case Unlock:
		s.WriteString(fmt.Sprintf("%d", Unlock))
	case EntryExitLock:
		s.WriteString(fmt.Sprintf("%d", EntryExitLock))
	case EntryLock:
		s.WriteString(fmt.Sprintf("%d", EntryLock))
	case ExitLock:
		s.WriteString(fmt.Sprintf("%d", ExitLock))
	default:
		s.WriteString("0")
	}

	// RegisterType
	switch opt.RegisterType {
	case InputRegister:
		s.WriteString(fmt.Sprintf("%d", InputRegister))
	case OutputRegister:
		s.WriteString(fmt.Sprintf("%d", OutputRegister))
	case InputOutputRegister:
		s.WriteString(fmt.Sprintf("%d", InputOutputRegister))
	case WithoutRegister:
		s.WriteString(fmt.Sprintf("%d", WithoutRegister))
	default:
		s.WriteString("0")
	}

	// ControlType
	switch opt.ControlType {
	case WithoutControl:
		s.WriteString(fmt.Sprintf("%d", WithoutControl))
	case InputControl:
		s.WriteString(fmt.Sprintf("%d", InputControl))
	case OutputControl:
		s.WriteString(fmt.Sprintf("%d", OutputControl))
	case InputOutputControl:
		s.WriteString(fmt.Sprintf("%d", InputOutputControl))
	default:
		s.WriteString("0")
	}

	// PictoMapType
	switch opt.PictoMapType {
	case OnePicto_OneGreen_OneRed:
		s.WriteString(fmt.Sprintf("%d", OnePicto_OneGreen_OneRed))
	case TwoPicto_OneGreen_OneRed:
		s.WriteString(fmt.Sprintf("%d", TwoPicto_OneGreen_OneRed))
	case OnePicto_TwoGreen_OneRed:
		s.WriteString(fmt.Sprintf("%d", OnePicto_TwoGreen_OneRed))
	default:
		s.WriteString("0")
	}

	// ViolationControlType
	switch opt.ViolationControlType {
	case SendAlarm:
		s.WriteString(fmt.Sprintf("%d", SendAlarm))
	case LockUntilReset:
		s.WriteString(fmt.Sprintf("%d", LockUntilReset))
	case BuzzerActiveUntilReset:
		s.WriteString(fmt.Sprintf("%d", BuzzerActiveUntilReset))
	case LockAndBuzzerActiveUntilReset:
		s.WriteString(fmt.Sprintf("%d", LockAndBuzzerActiveUntilReset))
	default:
		s.WriteString("0")
	}

	// TurnstileType
	switch opt.TurnstileType {
	case StandarHalfBodyWithOpticalSensorKitLogicLevelOne:
		s.WriteString(fmt.Sprintf("%d", StandarHalfBodyWithOpticalSensorKitLogicLevelOne))
	case StandarHalfBodyWithOpticalSensorKitLogicLevelZero:
		s.WriteString(fmt.Sprintf("%d", StandarHalfBodyWithOpticalSensorKitLogicLevelZero))
	case WolgateHalfBodyWithOpticalSensorKitLogicLevelOne:
		s.WriteString(fmt.Sprintf("%d", WolgateHalfBodyWithOpticalSensorKitLogicLevelOne))
	case WolgateHalfBodyWithOpticalSensorKitLogicLevelZero:
		s.WriteString(fmt.Sprintf("%d", WolgateHalfBodyWithOpticalSensorKitLogicLevelZero))
	case WolklugPlusHalfBodyWithOpticalSensorKitLogicLevelOne:
		s.WriteString(fmt.Sprintf("%d", WolklugPlusHalfBodyWithOpticalSensorKitLogicLevelOne))
	case WolklugPlusHalfBodyWithOpticalSensorKitLogicLevelZero:
		s.WriteString(fmt.Sprintf("%d", WolklugPlusHalfBodyWithOpticalSensorKitLogicLevelZero))
	default:
		s.WriteString("0")
	}

	// ReleaseControlType, BQCType, CheckType, CheckPercent
	// ... similar a los anteriores

	switch opt.ReleaseControlType {
	case ReleaseWithContinueSignal:
		s.WriteString(fmt.Sprintf("%d", ReleaseWithContinueSignal))
	case ReleaseWithOnePulse:
		s.WriteString(fmt.Sprintf("%d", ReleaseWithOnePulse))
	case ReleaseWithContinueSignalAndAutolock:
		s.WriteString(fmt.Sprintf("%d", ReleaseWithContinueSignalAndAutolock))
	default:
		s.WriteString("0")
	}

	switch opt.BQCType {
	case WithBQC:
		s.WriteString(fmt.Sprintf("%d", WithBQC))
	case WithoutBQC:
		s.WriteString(fmt.Sprintf("%d", WithoutBQC))
	default:
		s.WriteString("0")
	}

	switch opt.CheckType {
	case RandomCheckIsEntry:
		s.WriteString(fmt.Sprintf("%d", RandomCheckIsEntry))
	case RandomCheckIsExit:
		s.WriteString(fmt.Sprintf("%d", RandomCheckIsExit))
	default:
		s.WriteString("0")
	}

	switch opt.CheckPercent {
	case CheckPercentDisable:
		s.WriteString("000")
	case CheckPercent_100:
		s.WriteString("001")
	case CheckPercent_50:
		s.WriteString("002")
	case CheckPercent_33:
		s.WriteString("003")
	case CheckPercent_25:
		s.WriteString("004")
	case CheckPercent_20:
		s.WriteString("005")
	case CheckPercent_16, CheckPercent_14, CheckPercent_12, CheckPercent_11, CheckPercent_10:
		s.WriteString("0006")
	case CheckPercent_9, CheckPercent_8, CheckPercent_7, CheckPercent_6, CheckPercent_4, CheckPercent_3:
		s.WriteString("016")
	case CheckPercent_1, CheckPercent_05:
		s.WriteString("200")
	default:
		s.WriteString("0")
	}

	// TimeoutEntry
	s.WriteString(fmt.Sprintf("%03d", opt.TimeoutEntry))

	// MaxAccumEntry
	s.WriteString(fmt.Sprintf("%03d", opt.MaxAccumEntry))

	// TimeoutSelenoideControlCorredera
	s.WriteString(fmt.Sprintf("%03d", opt.TimeoutSelenoideControl))

	// ButtonsControl
	switch opt.OneButtonsControl {
	case OneButtonsAreButtons:
		s.WriteString(fmt.Sprintf("%d", OneButtonsAreButtons))
	case OneButtonsRandomCheckEntryForce:
		s.WriteString(fmt.Sprintf("%d", OneButtonsRandomCheckEntryForce))
	case OneButtonsEntryRelease:
		s.WriteString(fmt.Sprintf("%d", OneButtonsEntryRelease))
	case OneButtonsEntryExitRelease:
		s.WriteString(fmt.Sprintf("%d", OneButtonsEntryExitRelease))
	}

	// ButtonsControl
	switch opt.TwoButtonsControl {
	case OneButtonsAreButtons:
		s.WriteString(fmt.Sprintf("%d", OneButtonsAreButtons))
	case OneButtonsRandomCheckEntryForce:
		s.WriteString(fmt.Sprintf("%d", OneButtonsRandomCheckEntryForce))
	case OneButtonsEntryRelease:
		s.WriteString(fmt.Sprintf("%d", OneButtonsEntryRelease))
	case OneButtonsEntryExitRelease:
		s.WriteString(fmt.Sprintf("%d", OneButtonsEntryExitRelease))
	}

	// CorrederaType, PictoControlType
	// ... similar a los anteriores

	switch opt.CorrederaType {
	case ExitCorredera:
		s.WriteString("S")
	case EntryCorredera:
		s.WriteString("E")
	default:
		s.WriteString("S")
	}

	switch opt.PictoControlType {
	case PwaciiControl:
		s.WriteString("S")
	case SerialExternalControl:
		s.WriteString("X")
	default:
		s.WriteString("S")
	}

	return s.String()
}

type OptsFunc func(*Opts)

func DefaultsOptions() Opts {
	return Opts{
		Lock:                    Unlock,                                           // 0
		RegisterType:            InputOutputRegister,                              // 3
		ControlType:             InputControl,                                     // 1
		PictoMapType:            TwoPicto_OneGreen_OneRed,                         // 1
		ViolationControlType:    SendAlarm,                                        // 0
		TurnstileType:           WolgateHalfBodyWithOpticalSensorKitLogicLevelOne, // 4
		ReleaseControlType:      ReleaseWithOnePulse,                              // 1
		BQCType:                 WithoutBQC,                                       // 0
		CheckType:               RandomCheckIsExit,                                // 1
		CheckPercent:            CheckPercentDisable,                              // 000
		TimeoutEntry:            150,                                              // 150
		MaxAccumEntry:           0,                                                // 000
		TimeoutSelenoideControl: 0,                                                // 000
		OneButtonsControl:       OneButtonsAreButtons,                             // 1
		TwoButtonsControl:       TwoButtonsAreButtons,                             // 1
		CorrederaType:           ExitCorredera,                                    // S
		PictoControlType:        PwaciiControl,                                    // S
	}
}

func WithLock(lock Lock) OptsFunc {
	return func(opts *Opts) {
		opts.Lock = lock
	}
}

func WithRegisterType(registerType RegisterType) OptsFunc {
	return func(opts *Opts) {
		opts.RegisterType = registerType
	}
}

func WithControlType(controlType ControlType) OptsFunc {
	return func(opts *Opts) {
		opts.ControlType = controlType
	}
}

func WithPictoMapType(pictoMapType PictoMapType) OptsFunc {
	return func(opts *Opts) {
		opts.PictoMapType = pictoMapType
	}
}

func WithViolationControlType(violationControlType ViolationControlType) OptsFunc {
	return func(opts *Opts) {
		opts.ViolationControlType = violationControlType
	}
}

func WithTurnstileType(turnstileType TurnstileType) OptsFunc {
	return func(opts *Opts) {
		opts.TurnstileType = turnstileType
	}
}

func WithReleaseControlType(releaseControlType ReleaseControlType) OptsFunc {
	return func(opts *Opts) {
		opts.ReleaseControlType = releaseControlType
	}
}

func WithBQCType(bqcType BQCType) OptsFunc {
	return func(opts *Opts) {
		opts.BQCType = bqcType
	}
}

func WithCheckType(checkType CheckType) OptsFunc {
	return func(opts *Opts) {
		opts.CheckType = checkType
	}
}

func WithCheckPercent(checkPercent CheckPercent) OptsFunc {
	return func(opts *Opts) {
		opts.CheckPercent = checkPercent
	}
}

func WithOneButtonsControl(buttonsControl ButtonsControl) OptsFunc {
	return func(opts *Opts) {
		opts.OneButtonsControl = buttonsControl
	}
}

func WithTwoButtonsControl(buttonsControl ButtonsControl) OptsFunc {
	return func(opts *Opts) {
		opts.TwoButtonsControl = buttonsControl
	}
}

func WithCorrederaType(correderaType CorrederaType) OptsFunc {
	return func(opts *Opts) {
		opts.CorrederaType = correderaType
	}
}

func WithPictoControlType(pictoControlType PictoControlType) OptsFunc {
	return func(opts *Opts) {
		opts.PictoControlType = pictoControlType
	}
}

func WithTimeoutSelenoideControlCorredera(timeout TimeoutSelenoideControlCorredera) OptsFunc {
	return func(opts *Opts) {
		opts.TimeoutSelenoideControl = timeout
	}
}

func WithTimeoutEntry(timeout TimeoutEntry) OptsFunc {
	return func(opts *Opts) {
		opts.TimeoutEntry = timeout
	}
}

func WithMaxAccumEntry(maxAccum MaxAccumEntry) OptsFunc {
	return func(opts *Opts) {
		opts.MaxAccumEntry = maxAccum
	}
}

type Lock int

const (
	Unlock Lock = iota
	EntryLock
	ExitLock
	EntryExitLock
)

type RegisterType int

const (
	WithoutRegister RegisterType = iota
	InputRegister
	OutputRegister
	InputOutputRegister
)

type ControlType int

const (
	WithoutControl ControlType = iota
	InputControl
	OutputControl
	InputOutputControl
)

type PictoMapType int

const (
	OnePicto_OneGreen_OneRed PictoMapType = iota
	TwoPicto_OneGreen_OneRed
	OnePicto_TwoGreen_OneRed
)

type ViolationControlType int

const (
	SendAlarm ViolationControlType = iota
	LockUntilReset
	BuzzerActiveUntilReset
	LockAndBuzzerActiveUntilReset
)

type TurnstileType int

const (
	StandarHalfBodyWithOpticalSensorKitLogicLevelOne TurnstileType = iota
	StandarHalfBodyWithOpticalSensorKitLogicLevelZero
	WolklugPlusHalfBodyWithOpticalSensorKitLogicLevelOne
	WolklugPlusHalfBodyWithOpticalSensorKitLogicLevelZero
	WolgateHalfBodyWithOpticalSensorKitLogicLevelOne
	WolgateHalfBodyWithOpticalSensorKitLogicLevelZero
)

type ReleaseControlType int

const (
	ReleaseWithContinueSignal ReleaseControlType = iota
	ReleaseWithOnePulse
	ReleaseWithContinueSignalAndAutolock
)

type BQCType int

const (
	WithoutBQC BQCType = iota
	WithBQC
)

type CheckType int

const (
	RandomCheckIsEntry CheckType = iota
	RandomCheckIsExit
)

type CheckPercent int

const (
	CheckPercentDisable CheckPercent = iota
	CheckPercent_100
	CheckPercent_50
	CheckPercent_33
	CheckPercent_25
	CheckPercent_20
	CheckPercent_16
	CheckPercent_14
	CheckPercent_12
	CheckPercent_11
	CheckPercent_10
	CheckPercent_9
	CheckPercent_8
	CheckPercent_7
	CheckPercent_6
	CheckPercent_4
	CheckPercent_3
	CheckPercent_1
	CheckPercent_05
)

type TimeoutEntry uint8
type MaxAccumEntry uint8
type TimeoutSelenoideControlCorredera uint8

type ButtonsControl int

const (
	OneButtonsAreButtons            ButtonsControl = 1
	OneButtonsRandomCheckEntryForce ButtonsControl = 2
	OneButtonsEntryRelease          ButtonsControl = 4
	OneButtonsEntryExitRelease      ButtonsControl = 8
	TwoButtonsAreButtons            ButtonsControl = 1
	TwoButtonsRandomCheckEntryForce ButtonsControl = 2
	TwoButtonsEntryRelease          ButtonsControl = 4
)

type CorrederaType int

const (
	ExitCorredera CorrederaType = iota
	EntryCorredera
)

type PictoControlType int

const (
	PwaciiControl PictoControlType = iota
	SerialExternalControl
)
