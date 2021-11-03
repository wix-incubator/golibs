package logging

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

const appName = "test-app"
const logFolder = "."

func loadLogFile(fileName string) (retVal []map[string]interface{}) {
	bytesRead, _ := ioutil.ReadFile(fileName)
	file_content := string(bytesRead)
	lines := strings.Split(strings.Trim(file_content, " \n"), "\n")

	retVal = make([]map[string]interface{}, len(lines))
	for i, line := range lines {
		var result map[string]interface{}

		// Unmarshal or Decode the JSON to the interface.
		json.Unmarshal([]byte(line), &result)
		retVal[i] = result
	}
	return retVal
}
func Test_LevelConfig(t *testing.T) {
	SetLogConfig(Config{
		Level:  "debug",
		Colors: false,
		AdditionalFields: LoggerFields{
			Dc:              "42",
			ArtifactID:      "com.wixpress.artifact",
			ArtifactVersion: "1.0.1",
			Hostname:        "pod-1",
		},
	})

	entry := GetLog("test")

	assert.Equal(t, entry.Logger.Level, logrus.DebugLevel)
}

func Test_LogFileCreation(t *testing.T) {
	expectedFileName := path.Join(logFolder, fmt.Sprintf("%s_logstash_json.log", appName))
	os.Remove(expectedFileName)
	defer func() {
		os.Remove(expectedFileName)
	}()

	SetLogConfig(Config{
		LogsFolder:    ".",
		LogToJsonFile: true,
		AppName:       "test-app",
		Level:         "debug",
		Colors:        false,
		AdditionalFields: LoggerFields{
			Dc:              "42",
			ArtifactID:      "com.wixpress.artifact",
			ArtifactVersion: "1.0.1",
			Hostname:        "pod-1",
		},
	})

	logEntry := GetLog("test")
	logEntry.WithField("action", "someaction").Info("Test")
	etnries := loadLogFile(expectedFileName)
	assert.NotEmpty(t, etnries)
	assert.Len(t, etnries, 2)
}

func Test_LogDataConversion(t *testing.T) {
	expectedFileName := path.Join(logFolder, fmt.Sprintf("%s_logstash_json.log", appName))
	os.Remove(expectedFileName)
	defer func() {
		os.Remove(expectedFileName)
	}()

	SetLogConfig(Config{
		LogsFolder:    ".",
		LogToJsonFile: true,
		AppName:       "test-app",
		Level:         "debug",
		Colors:        false,
		AdditionalFields: LoggerFields{
			Dc:              "42",
			ArtifactID:      "com.wixpress.artifact",
			ArtifactVersion: "1.0.1",
			Hostname:        "pod-1",
		},
	})

	logEntry := GetLog("test")
	logEntry.WithField("action", "someaction").Info("Test")
	entries := loadLogFile(expectedFileName)
	assert.NotEmpty(t, entries)
	assert.Len(t, entries, 2)
	dataEntry := entries[1]["data"]
	datas := make(map[string]string)
	json.Unmarshal([]byte(dataEntry.(string)), &datas)
	assert.Equal(t, datas["action"], "someaction")
}
