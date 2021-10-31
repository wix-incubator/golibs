package logging

import (
	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	FileName   string
	Level      string
	Colors     bool
	Properties LogProperties
}

type LogProperties struct {
	DcName         string
	AppName        string
	PodName        string
	ArifactName    string
	ArifactVersion string
}

var loggerEntry *logrus.Entry
var colorSupport aurora.Aurora

func SetLogConfig(config Config) {
	configureLogger(config)
}

func GetCS() (retVal aurora.Aurora) {
	return colorSupport
}

func GetLog(obj string) (retVal *logrus.Entry) {
	retVal = loggerEntry.WithField("obj", obj)
	return retVal
}

func init() {
	configureLogger(Config{
		Level:      "debug",
		Colors:     false,
		Properties: LogProperties{},
	})
}

func configureLogger(config Config) {
	lock := sync.Mutex{}
	lock.Lock()
	defer lock.Unlock()

	levelStr := config.Level
	logger := logrus.New()
	logLevel, err := logrus.ParseLevel(levelStr)
	if err != nil {
		panic(err)
	}
	logger.SetLevel(logLevel)

	if config.FileName != "" {
		logFileName := config.FileName
		fileHook := NewJsonLogFileHook(logFileName, logLevel, config.Properties)
		logger.Hooks.Add(fileHook)
	}

	colorSupport = aurora.NewAurora(config.Colors)
	loggerEntry = logrus.NewEntry(logger)
	GetLog("logging").Info("Logging module configured successfully with", config)
}
