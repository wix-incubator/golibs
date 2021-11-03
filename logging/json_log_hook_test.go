package logging

import (
	"bytes"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"os"
	"path"
	"testing"
)

func Test_NewJsonLogFileHook(t *testing.T) {
	logFileName := path.Join(os.TempDir(), randomStr()+"-file.log")
	defer os.Remove(logFileName)
	loggerFields := LoggerFields{
		Dc:              "42",
		ArtifactID:      "com.wixpress.artifact",
		ArtifactVersion: "1.0.1",
		Hostname:        "pod-1",
	}
	obj := NewJsonLogFileHook(logFileName, loggerFields, logrus.TraceLevel)

	testNewJsonLogHook(obj, t)
}

func Test_NewJsonLogHook(t *testing.T) {
	loggerFields := LoggerFields{
		Dc:              "42",
		ArtifactID:      "com.wixpress.artifact",
		ArtifactVersion: "1.0.1",
		Hostname:        "pod-1",
	}
	obj := NewJsonLogHook(logrus.TraceLevel, loggerFields, new(bytes.Buffer))
	testNewJsonLogHook(obj, t)
}

func Test_JsonLogFireProperties(t *testing.T) {
	expect := assert.New(t)

	jsonMap := fireAndInterceptAsMapWith()
	expect.NotNil(jsonMap)
	expect.NotNil(jsonMap["data"])
	expect.NotNil(jsonMap["timestamop"])
}

func fireAndInterceptAsMapWith() map[string]string {
	buffer := new(bytes.Buffer)
	loggerFields := LoggerFields{
		Dc:              "42",
		ArtifactID:      "com.wixpress.artifact",
		ArtifactVersion: "1.0.1",
		Hostname:        "pod-1",
	}
	obj := NewJsonLogHook(logrus.TraceLevel, loggerFields, buffer)

	entry := newLogEntry(logrus.New(), loggerFields)
	entry.Level = logrus.TraceLevel
	obj.Fire(entry)

	jsonMap := make(map[string]string)
	json.Unmarshal([]byte(buffer.String()), &jsonMap)

	return jsonMap
}

func testNewJsonLogHook(hook *JsonLogHook, t *testing.T) {
	expect := assert.New(t)

	expect.NotNil(hook)
	expect.Equal(hook.levels, logrus.AllLevels)
}

func randomStr() string {
	return "somestring"
}
