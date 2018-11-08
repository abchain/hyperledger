package utils

import (
	google_protobuf "github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

// if t1 > t2, return 1
// else if t1 < t2, return -1
// esle (t1 == t2), return 0
func ComparePBTimestamp(t1 *google_protobuf.Timestamp, t2 *google_protobuf.Timestamp) int {

	if t1.GetSeconds() > t2.GetSeconds() {
		return 1
	} else if t1.GetSeconds() < t2.GetSeconds() {
		return -1
	} else if t1.GetNanos() > t2.GetNanos() {
		return 1
	} else if t1.GetNanos() < t2.GetNanos() {
		return -1
	} else {
		return 0
	}
}

func CreatePBTimestamp(t time.Time) *google_protobuf.Timestamp {

	ret := new(google_protobuf.Timestamp)

	ret.Seconds = t.Unix()
	ret.Nanos = int32(t.Nanosecond())

	return ret
}

func ConvertPBTimestamp(t *google_protobuf.Timestamp) time.Time {

	return time.Unix(t.GetSeconds(), int64(t.GetNanos()))
}
