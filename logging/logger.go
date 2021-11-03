package logging

import (
	"fmt"
	"path"
	"sync"

	"github.com/logrusorgru/aurora"
	"github.com/sirupsen/logrus"
)

type LoggerFields struct {
	ArtifactID      string
	ArtifactVersion string
	Hostname        string
	Dc              string
}

type Config struct {
	AppName          string
	LogsFolder       string
	LogToJsonFile    bool
	Level            string
	Colors           bool
	AdditionalFields LoggerFields
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
		Level:  "debug",
		Colors: false,
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

	if config.LogToJsonFile {
		shortLogFileName := fmt.Sprintf("%s_logstash_json.log", config.AppName)
		logFileName := path.Join(config.LogsFolder, shortLogFileName)
		fileHook := NewJsonLogFileHook(logFileName, config.AdditionalFields, logLevel)
		logger.Hooks.Add(fileHook)
	}

	colorSupport = aurora.NewAurora(config.Colors)
	loggerEntry = logrus.NewEntry(logger)
	GetLog("logging").Info("Logging module configured successfully with", config)
}
