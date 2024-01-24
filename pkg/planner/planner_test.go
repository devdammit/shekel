package planner_test

import (
	"testing"
	"time"

	"github.com/devdammit/shekel/pkg/planner"
	"github.com/devdammit/shekel/pkg/pointer"
	"github.com/stretchr/testify/assert"
)

var startAt = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
var endAt = time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

func TestPlanner_GetDates(t *testing.T) {
	t.Run("should return 12 days for daily planner with count 12", func(t *testing.T) {
		p := planner.NewPlanner(1, planner.PlanRepeatIntervalDaily, endAt, startAt, pointer.Ptr(uint32(12)), nil)

		dates, err := p.GetDates()

		assert.NoError(t, err)
		assert.Equal(t, 12, len(dates))
		assert.Equal(t, time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC), dates[1])
	})

	t.Run("should return 73 dates for daily planner with interval every 5 days", func(t *testing.T) {
		p := planner.NewPlanner(5, planner.PlanRepeatIntervalDaily, endAt, startAt, nil, nil)
		dates, err := p.GetDates()

		assert.NoError(t, err)
		assert.Equal(t, 73, len(dates))
		assert.Equal(t, time.Date(2024, 1, 6, 0, 0, 0, 0, time.UTC), dates[1])
	})

	t.Run("should return 52 dates for weekly planner with count 52", func(t *testing.T) {
		p := planner.NewPlanner(1, planner.PlanRepeatIntervalWeekly, endAt, startAt, pointer.Ptr(uint32(52)), []time.Weekday{time.Monday})
		dates, err := p.GetDates()

		assert.NoError(t, err)
		assert.Equal(t, 52, len(dates))
		assert.Equal(t, time.Date(2024, 1, 8, 0, 0, 0, 0, time.UTC), dates[1])
	})

	t.Run("should return 6 dates for weekly planner with interval 2 weeks", func(t *testing.T) {
		endAtDate := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
		p := planner.NewPlanner(2, planner.PlanRepeatIntervalWeekly, endAtDate, startAt, nil, []time.Weekday{time.Monday, time.Wednesday, time.Friday})

		dates, err := p.GetDates()

		assert.NoError(t, err)
		assert.Equal(t, 6, len(dates))
		assert.Equal(t, time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC), dates[1])
	})

	t.Run("should return 12 dates for monthly planner with count 12", func(t *testing.T) {
		p := planner.NewPlanner(1, planner.PlanRepeatIntervalMonthly, endAt, startAt, pointer.Ptr(uint32(12)), nil)

		dates, err := p.GetDates()

		assert.NoError(t, err)
		assert.Equal(t, 12, len(dates))
		assert.Equal(t, time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC), dates[1])
	})

	t.Run("should return 3 dates for yearly planner with count 3", func(t *testing.T) {
		endAtDate := time.Date(2030, 12, 31, 0, 0, 0, 0, time.UTC)

		p := planner.NewPlanner(1, planner.PlanRepeatIntervalYearly, endAtDate, startAt, pointer.Ptr(uint32(3)), nil)

		dates, err := p.GetDates()

		assert.NoError(t, err)
		assert.Equal(t, 3, len(dates))
		assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), dates[1])
	})
}
