package krabdb

import (
	"testing"
	"time"
)

func TestQuote(t *testing.T) {
	testCases := []struct {
		given    interface{}
		expected string
	}{
		{nil, "null"},
		{int64(420), "420"},
		{uint64(420), "420"},
		{float64(42.1), "42.1"},
		{true, "true"},
		{false, "false"},
		{[]byte{255, 128, 0}, `'\xff8000'`},
		{`krab`, `'krab'`},
		{`oh'krab`, `'oh''krab'`},
		{`oh\'krab`, `'oh\''krab'`},
		{
			time.Date(2020, time.March, 1, 23, 59, 59, 999999999, time.UTC),
			`'2020-03-01 23:59:59.999999Z'`,
		},
	}

	for i, tc := range testCases {
		actual := Quote(tc.given)

		if tc.expected != actual {
			t.Errorf("[%d] expected %s, but got %s", i, tc.expected, actual)
		}
	}
}
