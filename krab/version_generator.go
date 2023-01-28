package krab

import "time"

type VersionGenerator interface {
	Next() string
}

type TimestampVersionGenerator struct{}

func (g *TimestampVersionGenerator) Next() string {
	version := time.Now().UTC().Format("20060102_150405") // YYYYMMDD_HHMMSS
	return version
}
