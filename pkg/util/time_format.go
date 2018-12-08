package util

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
	"time"
)

func ConvertTimestampToTime(ts *timestamp.Timestamp) (time.Time, error) {
	newTime, err := ptypes.Timestamp(ts)
	if err != nil {
		return time.Now().UTC(), errors.Wrap(err, "failed to decode timestamp")
	}
	return newTime, nil
}

func ConvertTimeToTimestamp(time time.Time) (*timestamp.Timestamp, error) {
	ts, err := ptypes.TimestampProto(time)
	if err != nil {
		return nil, errors.Wrap(err, "failed to decode time")
	}
	return ts, nil
}
