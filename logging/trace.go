package logging

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const FieldNameAction = "action"
const FieldNameTraceId = "trace_id"
const FieldNameSegment = "segment"
const FieldNameMarker = "marker"
const FieldNameDuration = "duration_sec"
const MarkerStart = "start"
const MarkerEnd = "end"

type Trace interface {
	StartSegment(segmentName string, args ...interface{}) Segment
	NewSegment() SegmentBuilder
	AddField(name string, value interface{}) Trace
	Log() *logrus.Entry
	Id() string
}

type trace struct {
	logger *logrus.Entry
	name   string
	id     string
}

func NewTrace(action string, logger *logrus.Entry) Trace {
	id := uuid.New().String()

	return NewTraceWithId(id, action, logger)
}

func NewTraceWithId(id string, action string, logger *logrus.Entry) Trace {
	return &trace{
		logger: logger,
		name:   action,
		id:     id,
	}
}

func (t *trace) StartSegment(segmentName string, args ...interface{}) Segment {
	return t.NewSegment().Start(segmentName, args...)
}

func (t *trace) Log() *logrus.Entry {
	return baseEntryForTrace(t)
}

func (t *trace) NewSegment() SegmentBuilder {
	return &segmentBuilder{
		parent: t,
		logger: t.logger,
	}
}

func (t *trace) AddField(name string, value interface{}) Trace {
	t.logger = t.logger.WithField(name, value)
	return t
}

func (t *trace) Id() string {
	return t.id
}

func baseEntryForTrace(trace *trace) *logrus.Entry {
	return trace.logger.WithFields(
		logrus.Fields{
			FieldNameTraceId: trace.id,
			FieldNameAction:  trace.name,
		})
}
