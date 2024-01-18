package datetime_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
)

func TestMustParseDate(t *testing.T) {
	t.Run("should return correct date", func(t *testing.T) {
		d := datetime.MustParseDate("2019-03-29")
		year, month, day := d.Date()
		assert.Equal(t, 2019, year)
		assert.Equal(t, time.March, month)
		assert.Equal(t, 29, day)
	})
	t.Run("should panic on incorrect date", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = datetime.MustParseDate("2019-03-99")
		})
		assert.Panics(t, func() {
			_ = datetime.MustParseDate("2019-33-09")
		})
		assert.Panics(t, func() {
			_ = datetime.MustParseDate("2019-03-9")
		})
		assert.Panics(t, func() {
			_ = datetime.MustParseDate("2019-3-09")
		})
	})
}

func TestDate_MarshalJSON(t *testing.T) {
	date := datetime.MustParseDate("2019-03-09")
	data, err := json.Marshal(date)
	assert.NoError(t, err)
	assert.Equal(t, `"2019-03-09"`, string(data))
}

func TestDate_UnmarshalJSON(t *testing.T) {
	expect := datetime.MustParseDate("2019-03-09")
	var d datetime.Date
	err := json.Unmarshal([]byte(`"2019-03-09"`), &d)
	assert.NoError(t, err)
	assert.Equal(t, expect, d)
}
