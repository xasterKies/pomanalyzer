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

// Controls the interval timer
func tick(ctx context.Context, id int64, config *IntervalConfig,
		start, periodic, end Callback) error {
			
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()

			i, err := config.repo.ByID(id)
			if err != nil {
				
			}

		}