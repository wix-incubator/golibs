package logging

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewTraceObjectModel(t *testing.T) {
	_, entry := newTestLogger()

	trace := NewTrace(randomStr(), entry)

	segment := trace.StartSegment(randomStr())

	assert.NotNil(t, segment)
	assert.Equal(t, trace, segment.Parent())
}

func Test_TraceLogEntriesShouldContainTraceFields(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()

	trace := NewTrace(expectedAction, entry)

	trace.Log().Info()

	assertLastEntryWithAction(t, expectedAction, hook)
	assertLastEntryHasTraceId(t, hook)
}

func Test_TraceAddFieldShouldAddTheSpecifiedFieldToTheTraceLogger(t *testing.T) {
	hook, entry := newTestLogger()
	expectedFieldName := randomStr()
	expectedFieldValue := randomStr()

	trace := NewTrace(randomStr(), entry).
		AddField(expectedFieldName, expectedFieldValue)

	trace.Log().Info()
	assertLastEntryHasFieldWith(expectedFieldName, expectedFieldValue, hook, t)

	trace.NewSegment().Start(randomStr())
	assertLastEntryHasFieldWith(expectedFieldName, expectedFieldValue, hook, t)
}

func Test_SegmentStartAndEndFunctionsShouldProduceLogEntriesWithStartAndEndMarkersRespectively(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()

	trace := NewTrace(expectedAction, entry)

	segment := trace.StartSegment(expectedSegment)

	assertLastEntryWithStartMarkerAndWith(t, expectedAction, expectedSegment, hook)

	segment.End()
	assertLastEntryWithEndMarkerAndWith(t, expectedAction, expectedSegment, hook)
}

func Test_SegmentBuilderCreatedSegmentStartAndEndFunctionsShouldProduceLogEntriesWithStartAndEndMarkersRespectively(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()

	builder := NewTrace(expectedAction, entry).NewSegment()

	segment := builder.Start(expectedSegment)

	assertLastEntryWithStartMarkerAndWith(t, expectedAction, expectedSegment, hook)

	segment.End()
	assertLastEntryWithEndMarkerAndWith(t, expectedAction, expectedSegment, hook)
}

func Test_CustomSegmentFieldShouldBeIncludedInSegmentLogEntries(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedFieldName := randomStr()
	expectedFieldValue := randomStr()

	builder := NewTrace(expectedAction, entry).
		NewSegment().
		WithField(expectedFieldName, expectedFieldValue)

	segment := builder.Start(expectedSegment)

	assertLastEntryWithStartMarkerAndWith(t, expectedAction, expectedSegment, hook)
	assertLastEntryHasFieldWith(expectedFieldName, expectedFieldValue, hook, t)

	segment.End()
	assertLastEntryWithEndMarkerAndWith(t, expectedAction, expectedSegment, hook)
	assertLastEntryHasFieldWith(expectedFieldName, expectedFieldValue, hook, t)
}

func Test_AllCustomSegmentFieldsShouldBeIncludedInSegmentLogEntries(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedFields := map[string]interface{}{
		randomStr(): randomStr(),
		randomStr(): randomStr(),
	}

	builder := NewTrace(expectedAction, entry).
		NewSegment().
		WithFields(expectedFields)

	segment := builder.Start(expectedSegment)

	assertLastEntryWithStartMarkerAndWith(t, expectedAction, expectedSegment, hook)
	assertLastEntryHasAllFields(expectedFields, hook, t)

	segment.End()
	assertLastEntryWithEndMarkerAndWith(t, expectedAction, expectedSegment, hook)
	assertLastEntryHasAllFields(expectedFields, hook, t)
}

func Test_SegmentLogEtriesShouldIncludeSegmentFields(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()

	trace := NewTrace(expectedAction, entry)
	segment := trace.StartSegment(expectedSegment)
	segment.Log().Info()

	assertLastEntryWithAction(t, expectedAction, hook)
	assertLastEntryWithSegment(t, expectedSegment, hook)
	assertLastEntryHasTraceId(t, hook)
}

func Test_SegmentEndWithArgsShouldProduceAMessage(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()

	trace := NewTrace(expectedAction, entry)

	segment := trace.StartSegment(expectedSegment)

	segment.End("Error: ", errors.New("test error"))
	assert.Equal(t, "Error: test error", hook.LastEntry().Message)
}

func Test_SegmentEndWithErrorIfWithErrorShouldProduceErrorLogEntry(t *testing.T) {
	hook, entry := newTestLogger()

	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMessage := randomStr()

	trace := NewTrace(expectedAction, entry)

	segment := trace.StartSegment(expectedSegment)

	segment.EndWithErrorIf(errors.New(expectedMessage))

	assertLastEntryWithEndMarkerAnErrorAndWith(t,
		expectedAction,
		expectedSegment,
		logrus.ErrorLevel,
		expectedMessage,
		hook)
}

func Test_SegmentEndWithErrorIfWithNilShouldProduceInfoLogEntry(t *testing.T) {
	hook, entry := newTestLogger()

	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMessage := randomStr()

	trace := NewTrace(expectedAction, entry)

	segment := trace.StartSegment(expectedSegment)

	segment.EndWithErrorIf(nil, expectedMessage)

	assertLastEntryWithEndMarkerAnErrorAndWith(t,
		expectedAction,
		expectedSegment,
		logrus.InfoLevel,
		expectedMessage,
		hook)
}

func Test_SegmentEndWithWarningIfWithErrorShouldProduceWarningLogEntry(t *testing.T) {
	hook, entry := newTestLogger()

	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMessage := randomStr()

	trace := NewTrace(expectedAction, entry)
	segment := trace.StartSegment(expectedSegment)

	segment.EndWithWarningIf(errors.New(expectedMessage))

	assertLastEntryWithEndMarkerAnErrorAndWith(t,
		expectedAction,
		expectedSegment,
		logrus.WarnLevel,
		expectedMessage,
		hook)
}

func Test_SegmentEndWithWarningIfWithNilShouldProduceInfoLogEntry(t *testing.T) {
	hook, entry := newTestLogger()

	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMessage := randomStr()

	trace := NewTrace(expectedAction, entry)

	segment := trace.StartSegment(expectedSegment)

	segment.EndWithWarningIf(nil, expectedMessage)

	assertLastEntryWithEndMarkerAnErrorAndWith(t,
		expectedAction,
		expectedSegment,
		logrus.InfoLevel,
		expectedMessage,
		hook)
}

func Test_SegmentMarkShouldProduceLogEntryWithAMarkerField(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMarker := randomStr()

	trace := NewTrace(expectedAction, entry)

	segment := trace.StartSegment(expectedSegment)

	segment.Mark(expectedMarker)
	assertLastEntryWithMarkerAndWith(t, expectedAction, expectedSegment, expectedMarker, hook)
}

func Test_WithErrorMarkersOnlyShouldSkipNonErrorMarkerEvents(t *testing.T) {
	hook, entry := newTestLogger()

	segment := NewTrace(randomStr(), entry).
		NewSegment().
		WithErrorMarkersOnly().
		Start(randomStr())
	assert.Empty(t, hook.LastEntry())

	segment.Mark(randomStr())
	assert.Empty(t, hook.LastEntry())

	segment.EndWithWarningIf(errors.New(randomStr()))
	assert.Empty(t, hook.LastEntry())
}

func Test_WithErrorMarkersOnlyShouldProduceEndEventWithErrorWhenPresent(t *testing.T) {
	hook, entry := newTestLogger()

	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMessage := randomStr()

	segment := NewTrace(expectedAction, entry).
		NewSegment().
		WithErrorMarkersOnly().
		Start(expectedSegment)
	assert.Empty(t, hook.LastEntry())

	segment.Mark(randomStr())
	assert.Empty(t, hook.LastEntry())

	segment.EndWithErrorIf(errors.New(expectedMessage))
	assertLastEntryWithEndMarkerAnErrorAndWith(t,
		expectedAction,
		expectedSegment,
		logrus.ErrorLevel,
		expectedMessage,
		hook)
}

func Test_WithDebugMarkersShouldProduceStartAndEndEventWithDebugLevel(t *testing.T) {
	hook, entry := newTestLogger()
	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedCustomMarker := randomStr()

	segment := NewTrace(expectedAction, entry).
		NewSegment().
		WithDebugMarkers().
		Start(expectedSegment)
	assertLastEntryWithMarkerAndLevelWith(t, logrus.DebugLevel, expectedAction, expectedSegment, "start", hook)

	segment.Mark(expectedCustomMarker)
	assertLastEntryWithMarkerAndLevelWith(t, logrus.DebugLevel, expectedAction, expectedSegment, expectedCustomMarker, hook)

	segment.End()
	assertLastEntryWithMarkerAndLevelWith(t, logrus.DebugLevel, expectedAction, expectedSegment, "end", hook)
}

func Test_WithDebugMarkersShouldProduceEndEventWithErrorLevelWhenAnErrorWasReported(t *testing.T) {
	hook, entry := newTestLogger()

	expectedAction := randomStr()
	expectedSegment := randomStr()
	expectedMessage := randomStr()

	segment := NewTrace(expectedAction, entry).
		NewSegment().
		WithDebugMarkers().
		Start(expectedSegment)

	segment.EndWithErrorIf(errors.New(expectedMessage))
	assertLastEntryWithEndMarkerAnErrorAndWith(t,
		expectedAction,
		expectedSegment,
		logrus.ErrorLevel,
		expectedMessage,
		hook)
}

func Test_AddFieldShouldAddTheSpecifiedFieldToTheSegment(t *testing.T) {
	hook, entry := newTestLogger()
	expectedFieldName := randomStr()
	expectedFieldValue := randomStr()

	segment := NewTrace(randomStr(), entry).
		NewSegment().
		Start(randomStr())
	assertLastEntryDoesNotHaveField(expectedFieldName, hook, t)

	segment.Log().Info()
	assertLastEntryDoesNotHaveField(expectedFieldName, hook, t)

	segment.
		AddField(expectedFieldName, expectedFieldValue).
		Log().
		Info()
	assertLastEntryHasFieldWith(expectedFieldName, expectedFieldValue, hook, t)

	segment.End()
	assertLastEntryHasFieldWith(expectedFieldName, expectedFieldValue, hook, t)
}

func newTestLogger() (*test.Hook, *logrus.Entry) {
	nullLogger, hook := test.NewNullLogger()
	nullLogger.Level = logrus.DebugLevel

	entry := logrus.NewEntry(nullLogger)

	return hook, entry
}

func assertLastEntryWithStartMarkerAndWith(t *testing.T, expectedAction string, expectedSegment string, hook *test.Hook) {
	assertLastEntryWithMarkerAndWith(t, expectedAction, expectedSegment, "start", hook)
}

func assertLastEntryWithEndMarkerAndWith(t *testing.T, action string, segment string, hook *test.Hook) {
	assertLastEntryWithMarkerAndWith(t, action, segment, "end", hook)
	assert.IsType(t, float32(0), hook.LastEntry().Data["duration_sec"])
}

func assertLastEntryWithEndMarkerAnErrorAndWith(t *testing.T, expectedAction string, expectedSegment string, expectedLevel logrus.Level, expectedMessage string, hook *test.Hook) {
	assertLastEntryWithEndMarkerAndWith(t, expectedAction, expectedSegment, hook)
	assert.Equal(t, expectedLevel, hook.LastEntry().Level)
	assert.Equal(t, expectedMessage, hook.LastEntry().Message)
}

func assertLastEntryWithMarkerAndWith(t *testing.T, expectedAction string, expectedSegment string, expectedMarker string, hook *test.Hook) {
	assertLastEntryWithMarkerAndLevelWith(t, logrus.InfoLevel, expectedAction, expectedSegment, expectedMarker, hook)
}

func assertLastEntryWithMarkerAndLevelWith(t *testing.T, level logrus.Level, expectedAction string, expectedSegment string, expectedMarker string, hook *test.Hook) {
	assertLastEntryWithAction(t, expectedAction, hook)
	assertLastEntryWithSegment(t, expectedSegment, hook)
	assertLastEntryHasTraceId(t, hook)
	assertLastEntryHasFieldWith("marker", expectedMarker, hook, t)
	assert.IsType(t, level, hook.LastEntry().Level)
}

func assertLastEntryWithAction(t *testing.T, expectedAction string, hook *test.Hook) {
	assertLastEntryHasFieldWith("action", expectedAction, hook, t)
}

func assertLastEntryWithSegment(t *testing.T, expectedSegment string, hook *test.Hook) {
	assertLastEntryHasFieldWith("segment", expectedSegment, hook, t)
}

func assertLastEntryHasTraceId(t *testing.T, hook *test.Hook) {
	assert.NotEmpty(t, hook.LastEntry().Data["trace_id"])
}

func assertLastEntryHasFieldWith(name string, value interface{}, hook *test.Hook, t *testing.T) {
	assert.NotNil(t, hook.LastEntry())
	assert.Equal(t, value, hook.LastEntry().Data[name])
}

func assertLastEntryDoesNotHaveField(name string, hook *test.Hook, t *testing.T) {
	assert.NotNil(t, hook.LastEntry())
	assert.Nil(t, hook.LastEntry().Data[name])
}

func assertLastEntryHasAllFields(expectedFields map[string]interface{}, hook *test.Hook, t *testing.T) {
	for key, value := range expectedFields {
		assertLastEntryHasFieldWith(key, value, hook, t)
	}
}
