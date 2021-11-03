package logging

import (
	"io"

	"encoding/json"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogFieldNames string

const (
	ArtifactIDField      LogFieldNames = "artifact_id"
	ArtifactVersionField               = "artifact_version"
	DCField                            = "dc"
	HostnameField                      = "HOSTNAME"
)

type JsonLogHook struct {
	levels       []logrus.Level
	fileLogEntry *logrus.Entry
}

func NewJsonLogFileHook(fileName string, fields LoggerFields, levelToSet logrus.Level) (retVal *JsonLogHook) {
	fileLG := &lumberjack.Logger{
		Filename:   fileName,
		MaxSize:    100,
		MaxBackups: 1,
		MaxAge:     30,
		Compress:   true,
	}

	return NewJsonLogHook(levelToSet, fields, fileLG)
}

func NewJsonLogHook(levelToSet logrus.Level, fields LoggerFields, writer io.Writer) (retVal *JsonLogHook) {
	logrusLogger := logrus.New()
	logrusLogger.Level = levelToSet
	logrusLogger.Out = writer
	logrusLogger.Formatter = NewLogJsonFormatter()

	newFileLogEntry := newLogEntry(logrusLogger, fields)

	levels := make([]logrus.Level, 0)
	for _, nextLevel := range logrus.AllLevels {
		levels = append(levels, nextLevel)
		if int32(nextLevel) >= int32(levelToSet) {
			break
		}
	}

	retVal = &JsonLogHook{
		levels:       levels,
		fileLogEntry: newFileLogEntry,
	}
	return retVal
}

func makeDataField(data logrus.Fields) (retVal string) {
	asBytes, _ := json.Marshal(data)
	retVal = string(asBytes)
	return retVal
}

// Fire is required to implement Logrus hook
func (hook *JsonLogHook) Fire(entry *logrus.Entry) error {
	type printMethod func(args ...interface{})
	var funcToCallForPrint printMethod
	dataField := makeDataField(entry.Data)

	//entryTolog := hook.fileLogEntry.WithFields(entry.Data)
	// entryTolog := hook.fileLogEntry.
	// 	WithField(string(DCField),entry.Data[string(DCField)]).
	// 	WithField(string(ArtifactIDField),entry.Data[string(ArtifactIDField)]).
	// 	WithField(string(ArtifactVersionField),entry.Data[string(ArtifactVersionField)]).
	// 	WithField(string(HostnameField),entry.Data[string(HostnameField)]).
	// 	WithField("data", dataField)
	entryTolog := hook.fileLogEntry.WithField("data", dataField)

	switch entry.Level {
	case logrus.DebugLevel:
		funcToCallForPrint = entryTolog.Debug
	case logrus.InfoLevel:
		funcToCallForPrint = entryTolog.Info
	case logrus.WarnLevel:
		funcToCallForPrint = entryTolog.Warn
	case logrus.ErrorLevel:
		funcToCallForPrint = entryTolog.Error
	case logrus.FatalLevel:
		funcToCallForPrint = entryTolog.Fatal
	case logrus.PanicLevel:
		funcToCallForPrint = entryTolog.Panic
	case logrus.TraceLevel:
		funcToCallForPrint = entryTolog.Trace
	}
	funcToCallForPrint(entry.Message)
	return nil
}

// Levels Required for logrus hook implementation
func (hook *JsonLogHook) Levels() []logrus.Level {
	return hook.levels
}

func newLogEntry(logger *logrus.Logger, fields LoggerFields) *logrus.Entry {
	return logrus.
		NewEntry(logger).
		WithField(string(DCField), fields.Dc).
		WithField(string(ArtifactIDField), fields.ArtifactID).
		WithField(string(ArtifactVersionField), fields.ArtifactVersion).
		WithField(string(HostnameField), fields.Hostname)
}
