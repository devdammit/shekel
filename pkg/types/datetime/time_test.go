package datetime_test

import (
	"testing"

	"github.com/devdammit/shekel/pkg/encoding/json"
	"github.com/devdammit/shekel/pkg/types/datetime"
	"github.com/stretchr/testify/assert"
)

func TestMustParseAsTime(t *testing.T) {
	t.Run("should return correct time", func(t *testing.T) {
		d := datetime.MustParseAsTime("07:11:30")
		hh, mm, ss := d.Hour(), d.Minute(), d.Second()
		assert.Equal(t, 7, hh)
		assert.Equal(t, 11, mm)
		assert.Equal(t, 30, ss)
	})

	t.Run("should panic on incorrect time", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = datetime.MustParseAsTime("25:11:30")
		})
		assert.Panics(t, func() {
			_ = datetime.MustParseAsTime("11:-11:")
		})
		assert.Panics(t, func() {
			_ = datetime.MustParseAsTime("-1:11:30")
		})
		assert.Panics(t, func() {
			_ = datetime.MustParseAsTime("12:65:30")
		})
	})
}

func TestTime_MarshalJSON(t *testing.T) {
	date := datetime.MustParseAsTime("07:11:30")
	data, err := json.Marshal(date)
	assert.NoError(t, err)
	assert.Equal(t, `"07:11:30"`, string(data))
}

func TestTime_UnmarshalJSON(t *testing.T) {
	expect := datetime.MustParseAsTime("07:11:30")
	var d datetime.Time
	err := json.Unmarshal([]byte(`"07:11:30"`), &d)
	assert.NoError(t, err)
	assert.Equal(t, expect, d)
}
