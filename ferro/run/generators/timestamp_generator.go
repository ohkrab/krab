package generators

import "time"

type GeneratedTimestamp time.Time

func (t GeneratedTimestamp) String() string {
	return time.Time(t).Format("20060102_150405") // YYYYMMDD_HHMMSS
}

type TimestampGenerator interface {
	Next() GeneratedTimestamp
}

type TimestampVersionGenerator struct{}

func (g *TimestampVersionGenerator) Next() GeneratedTimestamp {
	return GeneratedTimestamp(time.Now().UTC())
}
