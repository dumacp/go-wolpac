package pwacii

type CommandType int

const (
	ConfRequest CommandType = iota
	StatusRequest
	GeneralReleaseEnable
	GeneralReleaseDisable
	OneEntryAllow
	OneExitAllow
	OneEntryCancel
	OneExitCancel
	DropArmActive
	StatusViolationRequest
	CheckActiveEntry
	CheckActiveExit
)

func (cmd CommandType) String() string {
	switch cmd {
	case ConfRequest:
		return "ConfRequest"
	case StatusRequest:
		return "StatusRequest"
	case GeneralReleaseEnable:
		return "GeneralReleaseEnable"
	case GeneralReleaseDisable:
		return "GeneralReleaseDisable"
	case OneEntryAllow:
		return "OneEntryAllow"
	case OneExitAllow:
		return "OneExitAllow"
	case OneEntryCancel:
		return "OneEntryCancel"
	case OneExitCancel:
		return "OneExitCancel"
	case DropArmActive:
		return "DropArmActive"
	case StatusViolationRequest:
		return "StatusViolationRequest"
	case CheckActiveEntry:
		return "CheckActiveEntry"
	case CheckActiveExit:
		return "CheckActiveExit"
	default:
		return "UnknownCommand"
	}
}

func (cmd CommandType) Code() string {
	switch cmd {
	case ConfRequest:
		return "CG"
	case StatusRequest:
		return "SS"
	case GeneralReleaseEnable:
		return "LG1"
	case GeneralReleaseDisable:
		return "LG0"
	case OneEntryAllow:
		return "LE"
	case OneExitAllow:
		return "LS"
	case OneEntryCancel:
		return "TE"
	case OneExitCancel:
		return "TS"
	case DropArmActive:
		return "BC"
	case StatusViolationRequest:
		return "VR"
	case CheckActiveEntry:
		return "RPE"
	case CheckActiveExit:
		return "RPS"
	default:
		return ""
	}
}

func (cmd CommandType) WithResponse() bool {
	switch cmd {
	case ConfRequest:
	case StatusRequest:
	default:
		return false
	}
	return true
}

type Command struct {
	Type EventType
	Data string
}

func (d *Device) Command(cmd CommandType, data string) (string, error) {

	return command(d, cmd, data)
}
