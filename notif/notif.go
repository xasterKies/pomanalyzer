package notif

import "runtime"

// serverity constants
const (
	SeverityLow = iota
	SeverityNormal
	SeverityUrgent
)

type Severity int

// Notification schema
type Notify struct {
	title    string
	message  string
	severity Severity
}

// Initialize a new notification
func New(title, message string, severity Severity) *Notify {
	return &Notify{
		title:    title,
		message:  message,
		severity: severity,
	}
}

// Determines and returns each notification severity string.
func (s Severity) String() string {
	sev := "low"

	switch s {
	case SeverityLow:
		sev = "low"
	case SeverityNormal:
		sev = "normal"
	case SeverityUrgent:
		sev = "critical"
	}

	if runtime.GOOS == "windows" {
		switch s {
		case SeverityLow:
			sev = "Info"
		case SeverityNormal:
			sev = "Warning"
		case SeverityUrgent:
			sev = "Error"
		}
	}

	return sev
}
