package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

const expectedTimestampFieldName = "timestamp"
const expectedMessageFieldName = "message"

func Test_LogFieldNames(t *testing.T) {
	expect := assert.New(t)
	formatter := NewLogJsonFormatter()

	actualTimestampFieldName := formatter.FieldMap[logrus.FieldKeyTime]
	expect.Equal(expectedTimestampFieldName, actualTimestampFieldName)

	actualMessageFieldName := formatter.FieldMap[logrus.FieldKeyMsg]
	expect.Equal(expectedMessageFieldName, actualMessageFieldName)
}
