package datetimeoffset

import "time"

type DateTimeOffset struct {
	time.Time
}

func Unix(sec int64, nsec int64) DateTimeOffset {
	return DateTimeOffset{
		time.Unix(sec, nsec),
	}
}

func Now() DateTimeOffset {
	n := time.Now()
	return DateTimeOffset{
		time.Unix(n.Unix(), int64(n.Nanosecond())),
	}
}
