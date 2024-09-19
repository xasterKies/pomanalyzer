package pomodoro

import (
	"context"
	"errors"
	"time"
)

// Category constants
const (
	CategoryPomodoro   = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak  = "LongBreak"
)

// State constants
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ById(id int64) (Interval, error)
	Last() (Interval, error)
	Breaks(n int) ([]Interval, error)
}

var (
	ErrNoIntervals        = errors.New("No intervals")
	ErrIntervalNotRunning = errors.New("Interval not running")
	ErrIntervalCompleted  = errors.New("Interval is completed or cancelled")
	ErrInvalidState       = errors.New("Invalid State")
	ErrInvalidID          = errors.New("Invalid ID")
)

type IntervalConfig struct {
	repo               Repository
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

func NewConfig(repo Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig {
	
	c := &IntervalConfig{
		repo: repo,
		PomodoroDuration: 25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration: 15 * time.Minute,
	}

	if pomodoro > 0 {
		c.PomodoroDuration = pomodoro
	}

	if shortBreak > 0 {
		c.ShortBreakDuration = shortBreak
	}

	if longBreak > 0 {
		c.LongBreakDuration = longBreak
	}

	return c
}

func nextCategory(r Repository) (string, error) {
	li, err := r.Last()

	if err != nil && err == ErrNoIntervals {
		return CategoryPomodoro, nil
	}
	if err != nil {
		return "", err
	}
	if li.Category == CategoryLongBreak || li.Category == CategoryShortBreak {
		return CategoryPomodoro, nil
	}

	lastBreaks, err := r.Breaks(3)
	if err != nil {
		return "", err
	}

	if len(lastBreaks) < 3 {
		return CategoryShortBreak, nil
	}

	for _, i := range lastBreaks {
		if i.Category == CategoryLongBreak {
			return CategoryShortBreak, nil
		}
	}

	return CategoryLongBreak, nil
}

type Callback func(Interval)

func tick(ctx context.Context, id int64, config *IntervalConfig, start, periodic, end Callback) error {

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	i, err := config.repo.ById(id)
	if err != nil {
		return err
	}

	expire := time.After(i.PlannedDuration - i.ActualDuration)

	start(i)

	for {
		select {
		case <-ticker.C:
			i, err := config.repo.ById(id)
			if err != nil {
				return err
			}

			if i.State == StatePaused {
				return nil
			}

			i.ActualDuration += time.Second
			if err := config.repo.Update(i); err != nil {
				return err
			}
			periodic(i)
		case <-expire:
			i, err := config.repo.ById(id)
			if err != nil {
				return err
			}
			i.State = StateDone
			end(i)
			return config.repo.Update(i)
		}

	}
}

func newInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	category, err := nextCategory(config.repo)
	if err != nil {
		return i, err
	}

	i.Category = category

	switch category {
	case CategoryPomodoro:
		i.PlannedDuration = config.PomodoroDuration
	case CategoryShortBreak:
		i.PlannedDuration = config.ShortBreakDuration
	case CategoryLongBreak:
		i.PlannedDuration = config.LongBreakDuration
	}

	if i.ID, err = config.repo.Create(i); err != nil {
		return i, err
	}

	return i, nil
}

func GetInterval(config *IntervalConfig) (Interval, error) {
	i := Interval{}
	var err error

	i, err = config.repo.Last()

	if err != nil && i.State != StateCancelled && i.State != StateDone {
		return i, nil
	}

	return newInterval(config)
}

