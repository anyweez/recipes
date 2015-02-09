package log

import (
	"fmt"
	logrus "github.com/Sirupsen/logrus"
	"os"
	"time"
)

type Fields map[string]interface{}

type LogEvent struct {
	Id           int64
	Name         string
	StickyFields Fields
}

var (
	STATUS_START    uint = 0
	STATUS_OK       uint = 1
	STATUS_COMPLETE uint = 2
	STATUS_WARNING  uint = 3
	STATUS_ERROR    uint = 4
	STATUS_FATAL    uint = 5
)

func New(name string, sticky Fields) LogEvent {
	le := LogEvent{}
	le.Id = time.Now().UnixNano()
	le.Name = name
	le.StickyFields = sticky

	le.Update(STATUS_START, name, nil)
	return le
}

func updateFields(sticky Fields, fields Fields, eventid int64) Fields {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	fields["_process"] = os.Args[0]
	fields["_eventid"] = eventid

	// Copy sticky fields over whenever they exist.
	if sticky != nil {
		for k := range sticky {
			fields[k] = sticky[k]
		}
	}

	return fields
}

func (e *LogEvent) Update(status uint, message string, fields Fields) {
	fields = updateFields(e.StickyFields, fields, e.Id)

	switch status {
	case STATUS_START:
		Info(fmt.Sprintf("[STATUS_START] %s", message), fields)
		break
	case STATUS_OK:
		Info(fmt.Sprintf("[STATUS_OK] %s", message), fields)
		break
	case STATUS_COMPLETE:
		Info(fmt.Sprintf("[STATUS_COMPLETE] %s", message), fields)
		break
	case STATUS_WARNING:
		Warn(fmt.Sprintf("[STATUS_WARNING] %s", message), fields)
		break
	case STATUS_ERROR:
		Error(fmt.Sprintf("[STATUS_ERROR] %s", message), fields)
		break
	case STATUS_FATAL:
		Fatal(fmt.Sprintf("[STATUS_FATAL] %s", message), fields)
		break

	}
}

func Info(message string, fields Fields) {
	if fields == nil {
		fields = make(map[string]interface{})
	}

	fields["_process"] = os.Args[0]

	logrus.WithFields(logrus.Fields(fields)).Info(message)
}

func Warn(message string, fields Fields) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["_process"] = os.Args[0]

	logrus.WithFields(logrus.Fields(fields)).Warn(message)
}

func Error(message string, fields Fields) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["_process"] = os.Args[0]

	logrus.WithFields(logrus.Fields(fields)).Error(message)
}

func Fatal(message string, fields Fields) {
	if fields == nil {
		fields = make(map[string]interface{})
	}
	fields["_process"] = os.Args[0]

	logrus.WithFields(logrus.Fields(fields)).Fatal(message)
}
