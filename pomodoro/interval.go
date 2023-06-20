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
	StateNotStarted = iota // (iota) automatically increases the number of each line
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

// Representing pomodoro interval
type Interval struct {
	ID              int64
	StartTime       time.Time
	PlannedDuration time.Duration
	ActualDuration  time.Duration
	Category        string
	State           int
}

// Defining repository pattern for storing data
type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ByID(id int64) (Interval, error)
	Last() (Interval, error)
	Breaks(n int) ([]Interval, error)
}

// Representing particular errors app may return
var (
	ErrNoIntervals        = errors.New("No intervals")
	ErrIntervalNotRunning = errors.New("Interval not running")
	ErrIntervalCompleted  = errors.New("Interval is completed or cancelled")
	ErrInvalidState       = errors.New("Invalid State")
	ErrInvalidID          = errors.New("Invalid ID")
)

// Representing the configuration required to instantiate an interval
type IntervalConfig struct {
	repo               Repository
	PomodoroDuration   time.Duration
	ShortBreakDuration time.Duration
	LongBreakDuration  time.Duration
}

// Instantiating a new intervalConfig - Uses the values provided byt the user or
// set default values for each interval
func NewConfig(repo Repository, pomodoro, shortBreak, longBreak time.Duration) *IntervalConfig {

	c := &IntervalConfig{
		repo:               repo,
		PomodoroDuration:   25 * time.Minute,
		ShortBreakDuration: 5 * time.Minute,
		LongBreakDuration:  15 * time.Minute,
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

// This function retrieves the last interval from the repository and determines
// the next interval category based on the Pomodoro Technique rules. After each
// Pomodoro interval, there’s a short break, and after four Pomodoros, there’s
// a long break.
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

// Controls interval timer - This function uses the time.Ticker type and a loop to execute actions every
// second while the interval time progresses. It uses a select statement to take
// actions, executing periodically when the time.Ticker goes off
func tick(ctx context.Context, id int64, config *IntervalConfig,
	start, periodic, end Callback) error {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		i, err := config.repo.ByID(id)
		if err != nil {
			return err
		}
		expire := time.After(i.PlannedDuration - i.ActualDuration)

		start(i)

		for {
			select {
			case <-ticker.C:
				i, err := config.repo.ByID(id)
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
				i, err := config.repo.ByID(id)
				if err != nil {
					return err
				}
				i.State = StateDone
				end(i)
				return config.repo.Update(i)
			case <-ctx.Done():
				i, err := config.repo.ByID(id)
				if err != nil {
					return err
				}
				i.State = StateCancelled
				return config.repo.Update(i)
		}
	}
	}
