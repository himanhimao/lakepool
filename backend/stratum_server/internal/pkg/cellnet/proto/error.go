package proto

type ErrType int32

const (
	ErrOtherUnknown       ErrType = 20
	ErrJobNotFound        ErrType = 21
	ErrDuplicateShare     ErrType = 22
	ErrLowDifficultyShare ErrType = 23
	ErrUnauthorizedWorker ErrType = 24
	ErrNotSubscribed      ErrType = 25
)

func (p ErrType) ErrMsg() string {
	switch p {
	case ErrOtherUnknown:
		return "Other/Unknown"
	case ErrJobNotFound:
		return "Job not found"
	case ErrDuplicateShare:
		return "Duplicate Share"
	case ErrLowDifficultyShare:
		return "Low difficulty share"
	case ErrUnauthorizedWorker:
		return "Unauthorized worker"
	case ErrNotSubscribed:
		return "Not subscribed"
	default:
		return "Unknown"
	}
}

func (p ErrType) string() string {
	return p.ErrMsg()
}
