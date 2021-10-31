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

const expectedFieldNameDc = "dc"
const expectedFieldNameServiceName = "service"
const expectedFieldNameAppName = "app_name"
const expectedFieldNamePodName = "pod_name"
const expectedFieldNameInstance = "instance"

func Test_NewJsonLogFileHook(t *testing.T) {
	logFileName := path.Join(os.TempDir(), randomStr()+"-file.log")
	defer os.Remove(logFileName)

	obj := NewJsonLogFileHook(logFileName, logrus.TraceLevel, LogProperties{})

	testNewJsonLogHook(obj, t)
}

func Test_NewJsonLogHook(t *testing.T) {
	obj := NewJsonLogHook(logrus.TraceLevel, LogProperties{}, new(bytes.Buffer))
	testNewJsonLogHook(obj, t)
}

func Test_JsonLogFireProperties(t *testing.T) {
	expect := assert.New(t)

	var expectedProperties = &LogProperties{
		DcName:       randomStr(),
		AppName:      randomStr(),
		PodName:      randomStr(),
		ServiceName:  randomStr(),
		InstanceName: randomStr(),
	}

	jsonMap := fireAndInterceptAsMapWith(expectedProperties)

	expect.Equal(expectedProperties.DcName, jsonMap[expectedFieldNameDc])
	expect.Equal(expectedProperties.ServiceName, jsonMap[expectedFieldNameServiceName])
	expect.Equal(expectedProperties.AppName, jsonMap[expectedFieldNameAppName])
	expect.Equal(expectedProperties.PodName, jsonMap[expectedFieldNamePodName])
	expect.Equal(expectedProperties.InstanceName, jsonMap[expectedFieldNameInstance])
}

func Test_JsonLogFirePartialProperties(t *testing.T) {
	expect := assert.New(t)

	var expectedProperties = &LogProperties{
		DcName:  randomStr(),
		AppName: randomStr(),
	}

	jsonMap := fireAndInterceptAsMapWith(expectedProperties)

	expect.Equal(expectedProperties.DcName, jsonMap[expectedFieldNameDc])
	expect.Equal(expectedProperties.AppName, jsonMap[expectedFieldNameAppName])
	expect.Empty(jsonMap[expectedFieldNameServiceName])
	expect.Empty(jsonMap[expectedFieldNamePodName])
}

func fireAndInterceptAsMapWith(expectedProperties *LogProperties) map[string]string {
	buffer := new(bytes.Buffer)
	obj := NewJsonLogHook(logrus.TraceLevel, *expectedProperties, buffer)

	entry := newLogEntry(logrus.New(), expectedProperties)
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