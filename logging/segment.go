package logging

import (
	"github.com/sirupsen/logrus"
	"reflect"
	"time"
)

type Segment interface {
	Parent() Trace
	End(args ...interface{})
	EndWithErrorIf(err error, elseArgs ...interface{})
	EndWithWarningIf(err error, elseArgs ...interface{})
	Mark(marker string, args ...interface{}) Segment
	AddField(name string, value interface{}) Segment
	Log() *logrus.Entry

	start(args ...interface{})
}

type segment struct {
	logger          *logrus.Entry
	parent          *trace
	name            string
	startTime       time.Time
	markerLogMethod string
}

type errorMarkersOnlySegment struct {
	delegate Segment
}

func (s *segment) Parent() Trace {
	return s.parent
}

func (s *segment) End(args ...interface{}) {
	logMarkerEntry(s.endEntry(), s.markerLogMethod, args...)
}

func (s *segment) EndWithErrorIf(err error, elseArgs ...interface{}) {
	entry := s.endEntry()

	if err != nil {
		entry.Error(err)
	} else {
		logMarkerEntry(entry, s.markerLogMethod, elseArgs...)
	}
}

func (s *segment) EndWithWarningIf(err error, elseArgs ...interface{}) {
	entry := s.endEntry()

	if err != nil {
		entry.Warn(err)
	} else {
		logMarkerEntry(entry, s.markerLogMethod, elseArgs...)

	}
}

func (s *segment) Mark(marker string, args ...interface{}) Segment {

	entry := s.logger.WithFields(
		logrus.Fields{
			FieldNameSegment: s.name,
			FieldNameMarker:  marker,
		})

	logMarkerEntry(entry, s.markerLogMethod, args...)

	return s
}

func (s *segment) Log() *logrus.Entry {
	return s.logger.WithField(FieldNameSegment, s.name)
}

func (s *segment) AddField(name string, value interface{}) Segment {
	s.logger = s.logger.WithField(name, value)

	return s
}

func (s *segment) start(args ...interface{}) {
	entry := s.logger.WithField(FieldNameMarker, MarkerStart)

	logMarkerEntry(entry, s.markerLogMethod, args...)
}

func (s *segment) endEntry() *logrus.Entry {
	return s.logger.
		WithFields(
			logrus.Fields{
				FieldNameSegment:  s.name,
				FieldNameMarker:   MarkerEnd,
				FieldNameDuration: elapsedSec(s.startTime),
			})
}

func (s *errorMarkersOnlySegment) Parent() Trace {
	return s.delegate.Parent()
}

func (s *errorMarkersOnlySegment) End(args ...interface{}) {}

func (s *errorMarkersOnlySegment) EndWithErrorIf(err error, args ...interface{}) {
	s.delegate.EndWithErrorIf(err)
}

func (s *errorMarkersOnlySegment) EndWithWarningIf(err error, args ...interface{}) {}

func (s *errorMarkersOnlySegment) Mark(marker string, args ...interface{}) Segment {
	return s
}

func (s *errorMarkersOnlySegment) Log() *logrus.Entry {
	return s.delegate.Log()
}

func (s *errorMarkersOnlySegment) AddField(name string, value interface{}) Segment {
	s.delegate.AddField(name, value)

	return s
}

func (s *errorMarkersOnlySegment) start(args ...interface{}) {}

func elapsedSec(startTime time.Time) float32 {
	return float32(time.Since(startTime).Seconds())
}

func logMarkerEntry(entry *logrus.Entry, markerLogMethod string, args ...interface{}) {
	parameters := make([]reflect.Value, len(args))
	for i := range args {
		parameters[i] = reflect.ValueOf(args[i])
	}
	reflect.ValueOf(entry).MethodByName(markerLogMethod).Call(parameters)
}
