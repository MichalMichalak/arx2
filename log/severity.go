package log

type Severity string

const (
	SeverityError   Severity = "ERROR"
	SeverityWarning Severity = "WARNING"
	SeverityInfo    Severity = "INFO"
	SeverityDebug   Severity = "DEBUG"
)

func (s Severity) ThisAndAbove() map[Severity]struct{} {
	switch s {
	case SeverityError:
		return map[Severity]struct{}{SeverityError: {}}
	case SeverityWarning:
		return map[Severity]struct{}{SeverityError: {}, SeverityWarning: {}}
	case SeverityInfo:
		return map[Severity]struct{}{SeverityError: {}, SeverityWarning: {}, SeverityInfo: {}}
	default:
		return map[Severity]struct{}{SeverityError: {}, SeverityWarning: {}, SeverityInfo: {}, SeverityDebug: {}}
	}
}
