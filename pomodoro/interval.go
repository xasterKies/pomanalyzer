package pomodoro

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// Category constants
const (
	CategoryPomodoro = "Pomodoro"
	CategoryShortBreak = "ShortBreak"
	CategoryLongBreak = "LongBreak"
)

// State constants
const (
	StateNotStarted = iota
	StateRunning
	StatePaused
	StateDone
	StateCancelled
)

// Storing data in a repository pattern
type Repository interface {
	Create(i Interval) (int64, error)
	Update(i Interval) error
	ById(id int64) (Interval, error)
	Last() (Interval, errorr)
	Breaks(int) ([]Interval, error)
}

// Representing the particular errors business logic may return
var (
	ErrNoIntervals = errors.New("No intervals")
	ErrIntervalNotRunning = errors.New("Interval not running")
	ErrIntervalCompleted = errors.New("Interval is completed or cancelled")
	ErrInvalidState = errors.New("Invalid State")
	ErrInvalidID = errors.New("Invalid ID")
)

// Representing configuration required to instantiate an interval
type IntervalConfig struct {
	repo									Repository
	PomodoroDuration			time.Duration
	ShortBreakDuration		time.Duration
	LongBreakDuration			time.Duration
}

// createing a new Interval config struct
func Newconfig(repo Repository, pomodoro, shortBreak, longBreak
	time.Duration) *IntervalConfig {
	
		c := &IntervalConfig{
			repo:								repo,
			PomodoroDuration: 	25 * time.Minute,
			ShortBreakDuration: 5 * time.Minute,
			LongBreakDuration:	15 * time.Minute,
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


// =============== Methods for main interval type ===============

/* 
	* Retrives the last interval from the repository and determines the next interval
	*	based on the Pomodoro Technique rules
*/
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

	if len(lastBreaks) > 3 {
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

/*
	* Controls the interval timer.
	* Uses time.Ticker type and loop to execute actions every second while the interval time progress
*/
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
						if err !=  nil {
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
	
	/*
	 * Creates a new interval instance with the appropriate category and values
	*/
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

	// ================ API Functions - GetInterval(), Start(), Pause() =========================

	/*
	 * Retrives the last interval of the repository, returning if its active or returning an error
	*/
	func GetInterval(config *IntervalConfig) (Interval, error) {
		i := Interval{}
		var err error

		i, err = config.repo.Last()

		if err != nil && err != ErrNoIntervals {
			return i, err
		}

		if err == nil && i.State != StateCancelled && i.State != StateDone {
			return i, nil
		}

		return newInterval(config)
	}

	/*
	 * Checks the state of the current interval setting the appropriate options then calls
	 * tick() to time interval
	*/
	func (i Interval) Start(ctx content.Context, config *IntervalConfig,
			start, periodic, end Callback) error {
				
				switch i.State {
					case StateRunning:
						return nil
					case StateNotStarted:
						i.StartTime = time.Now()
						fallthrough
					case StatePaused:
						i.State = StateRunning
						if err := config.repo.Update(i); err != nil {
							return err
						}
						return tick(ctx, i.ID, config, start, periodic, end)
					case StateCancelled, StateDone:
						return fmt.Errorf("%w: Cannot start", ErrIntervalCompleted)
					default:
						return fmt.Errorf("%w: %d", ErrInvalidState, i.State)
				}
			}
	
	/*
	 * Verifies whether the instance of interval is running and pauses by settin the
	 * state to StatePaused
	*/
	func (i Interval) Pause(config *IntervalConfig) error {
		if i.State != StateRunning {
			return ErrIntervalNotRunning
		}

		i.State = StatePaused

		return config.repo.Update(i)
	}