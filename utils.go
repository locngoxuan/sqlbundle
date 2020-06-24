package sqlbundle

import "time"

const TIME_FORMAT = "20060102150405"

func makeTimeSequence() string {
	t := time.Now()
	return t.Format(TIME_FORMAT)
}
