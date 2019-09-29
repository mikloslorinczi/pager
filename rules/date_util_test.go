package rules

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var UTC = time.UTC

func Test_IsWorkHour_Exception(t *testing.T) {
	tm := time.Unix(1546357305, 0).In(UTC)
	require.False(t, IsWorkHour(tm, UTC), "2019 01 01 Was a holyday, IsWorkHour should return false")
}

func Test_IsWorkHour_Saturday(t *testing.T) {
	tm := time.Unix(1554547305, 0).In(UTC)
	require.False(t, IsWorkHour(tm, UTC), "2019 04 06 Was a Saturday, IsWorkHour should return false")
}

func Test_IsWorkHour_Monday(t *testing.T) {
	tm := time.Unix(1577119305, 0).In(UTC)
	require.True(t, IsWorkHour(tm, UTC), "2019 12 23 Was a Monday, IsWorkHour should return true")
}

func Test_IsWorkHour_Monday_OffWork(t *testing.T) {
	tm := time.Unix(1577120400, 0).In(UTC)
	require.False(t, IsWorkHour(tm, UTC), "at 17:00, IsWorkHour should return false")
}
