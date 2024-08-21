package logIPFS

import "fmt"

func (l *logger) info(msg interface{}) {
	l.System.Info(formatMessage(msg))
}

func (l *logger) warn(msg interface{}) {
	l.System.Warn(formatMessage(msg))
}

func (l *logger) error(msg interface{}) {
	l.System.Error(formatMessage(msg))
}

func formatMessage(msg interface{}) string {
	switch v := msg.(type) {
	case string:
		return v
	case error:
		return v.Error()
	default:
		return fmt.Sprintf("%v", v)
	}
}
