package datetime_test

import (
	"testing"
	"time"

	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
)

func TestDateTime(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2020-12-22T10:45:00+07:00")
	dt := datetime.DateTime{Time: ts}

	t.Run("should serialize with minute precision", func(t *testing.T) {
		assert.Equal(t, "2020-12-22 10:45", dt.String())
	})
}

func TestDateTime_LocalTimestamp(t *testing.T) {
	ts, _ := time.Parse(time.RFC3339, "2020-12-22T10:45:00+07:00")
	dt := datetime.DateTime{Time: ts}

	assert.NotEqual(t, dt.Unix(), dt.LocalTimestamp())
	assert.Equal(t, datetime.MustParseDateTime(dt.String()).Unix(), dt.LocalTimestamp())
}

func TestDateTimeMinutesSinceDateMidnight(t *testing.T) {
	date1 := datetime.MustParseDateTime("2020-12-22 10:45")
	date2 := datetime.MustParseDateTime("2020-12-23 00:00")

	assert.Equal(t, 1440, date2.MinutesSinceDateMidnight(date1))
}
