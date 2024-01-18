package planner

import (
	"errors"
	"time"
)

type PlanRepeatInterval string

const (
	PlanRepeatIntervalDaily   PlanRepeatInterval = "daily"
	PlanRepeatIntervalWeekly  PlanRepeatInterval = "weekly"
	PlanRepeatIntervalMonthly PlanRepeatInterval = "monthly"
	PlanRepeatIntervalYearly  PlanRepeatInterval = "yearly"
)

type Planner struct {
	intervalCount  uint32
	repeatInterval PlanRepeatInterval
	startAt        time.Time
	endAt          time.Time
	weekdays       []time.Weekday
	count          *uint32
}

func NewPlanner(intervalCount uint32, repeatInterval PlanRepeatInterval, endDate time.Time, startDate time.Time, count *uint32, weekdays []time.Weekday) *Planner {
	return &Planner{
		intervalCount:  intervalCount,
		repeatInterval: repeatInterval,
		startAt:        startDate,
		endAt:          endDate,
		weekdays:       weekdays,
		count:          count,
	}
}

func (p *Planner) GetDates() ([]time.Time, error) {
	if p.endAt.Before(p.startAt) {
		return nil, errors.New("end date is before start date")
	}

	switch p.repeatInterval {
	case PlanRepeatIntervalDaily:
		if p.weekdays != nil {
			return nil, errors.New("days of week cannot be set for daily interval")
		}

		dates, err := p.getDailyIntervalDates()
		if err != nil {
			return nil, err
		}

		return dates, nil

	case PlanRepeatIntervalWeekly:
		if p.weekdays == nil {
			return nil, errors.New("days of week cannot be empty for weekly interval")
		}

		dates, err := p.getWeeklyIntervalDates()
		if err != nil {
			return nil, err
		}

		return dates, nil

	case PlanRepeatIntervalMonthly:
		if p.weekdays != nil {
			return nil, errors.New("days of week cannot be set for monthly interval")
		}

		dates, err := p.getMonthlyIntervalDates()
		if err != nil {
			return nil, err
		}

		return dates, nil

	case PlanRepeatIntervalYearly:
		if p.weekdays != nil {
			return nil, errors.New("days of week cannot be set for yearly interval")
		}

		dates, err := p.getYearlyIntervalDates()
		if err != nil {
			return nil, err
		}

		return dates, nil
	}

	return nil, nil
}

func (p *Planner) getDailyIntervalDates() ([]time.Time, error) {
	var dates []time.Time

	if p.intervalCount == 0 {
		return nil, errors.New("interval count is 0")
	}

	for d := p.startAt; d.Before(p.endAt); d = d.AddDate(0, 0, int(p.intervalCount)) {
		if p.count != nil && len(dates) >= int(*p.count) {
			break
		}

		dates = append(dates, d)
	}

	return dates, nil
}

func (p *Planner) getWeeklyIntervalDates() ([]time.Time, error) {
	var dates []time.Time

	if p.intervalCount == 0 {
		return nil, errors.New("interval count is 0")
	}

	skip := 0

	for d := p.startAt; d.Before(p.endAt); d = d.AddDate(0, 0, 1) {
		if p.count != nil && len(dates) >= int(*p.count) {
			break
		}

		if skip > 0 {
			skip--
			continue
		}

		for _, dayOfWeek := range p.weekdays {
			if d.Weekday() == dayOfWeek {
				dates = append(dates, d)
			}
		}

		if d.Weekday() == time.Sunday {
			if p.intervalCount == 1 {
				skip = 0
			} else {
				skip = int(p.intervalCount) * 7
			}
		}
	}

	return dates, nil
}

func (p *Planner) getMonthlyIntervalDates() ([]time.Time, error) {
	var dates []time.Time

	if p.intervalCount == 0 {
		return nil, errors.New("interval count is 0")
	}

	for d := p.startAt; d.Before(p.endAt); d = d.AddDate(0, int(p.intervalCount), 0) {
		if p.count != nil && len(dates) >= int(*p.count) {
			break
		}

		dates = append(dates, d)
	}

	return dates, nil
}

func (p *Planner) getYearlyIntervalDates() ([]time.Time, error) {
	var dates []time.Time

	if p.intervalCount == 0 {
		return nil, errors.New("interval count is 0")
	}

	for d := p.startAt; d.Before(p.endAt); d = d.AddDate(int(p.intervalCount), 0, 0) {
		if p.count != nil && len(dates) >= int(*p.count) {
			break
		}

		dates = append(dates, d)
	}

	return dates, nil
}
