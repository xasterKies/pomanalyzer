package pomodoro_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/xasterKies/pomanalyzer/pomodoro"
)

func TestNewConfig(t *testing.T) {
	testCases := []struct {
		name   string
		input  [3]time.Duration
		expect pomodoro.IntervalConfig
	}{
		{
			name: "Default",
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   25 * time.Minute,
				ShortBreakDuration: 5 * time.Minute,
				LongBreakDuration:  15 * time.Minute,
			},
		},
		{
			name: "SingleInput",
			input: [3]time.Duration{
				20 * time.Minute,
			},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   20 * time.Minute,
				ShortBreakDuration: 5 * time.Minute,
				LongBreakDuration:  15 * time.Minute,
			},
		},
		{
			name: "MultiInput",
			input: [3]time.Duration{
				20 * time.Minute,
				10 * time.Minute,
				12 * time.Minute,
			},
			expect: pomodoro.IntervalConfig{
				PomodoroDuration:   20 * time.Minute,
				ShortBreakDuration: 10 * time.Minute,
				LongBreakDuration:  12 * time.Minute,
			},
		},
	}

	// Execute tests for NewConfig
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var repo pomodoro.Repository
			config := pomodoro.NewConfig(
				repo,
				tc.input[0],
				tc.input[1],
				tc.input[2],
			)

			if config.PomodoroDuration != tc.expect.PomodoroDuration {
				t.Errorf("Expected Pomodoro Duration %q, got %q instead\n",
					tc.expect.PomodoroDuration, config.PomodoroDuration)
			}
			if config.ShortBreakDuration != tc.expect.ShortBreakDuration {
				t.Errorf("Expected ShortBreak Duration %q, got %q instead\n",
					tc.expect.ShortBreakDuration, config.ShortBreakDuration)
			}
			if config.LongBreakDuration != tc.expect.LongBreakDuration {
				t.Errorf("Expected LongBreak Duration %q, got %q instead\n",
					tc.expect.LongBreakDuration, config.LongBreakDuration)
			}
		})
	}
}

func TestGetInterval(t *testing.T) {
	repo, cleanup := getRepo(t)
	defer cleanup()

	const duration = 1 * time.Millisecond
	config := pomodoro.NewConfig(repo, 3*duration, duration, 2*duration)

	for i := 1; i <= 16; i++ {
		var (
			expCategory string
			expDuration time.Duration
		)

		switch {
		case i%2 != 0:
			expCategory = pomodoro.CategoryPomodoro
			expDuration = 3 * duration
		case i%8 == 0:
			expCategory = pomodoro.CategoryLongBreak
			expDuration = 2 * duration
		case i%2 == 0:
			expCategory = pomodoro.CategoryShortBreak
			expDuration = duration
		}

		testName := fmt.Sprintf("%s: %d", expCategory, i)
		t.Run(testName, func(t *testing.T) {
			res, err := pomodoro.GetInterval(config)

			if err != nil {
				t.Errorf("Expected no error, got %q.\n", err)
			}

			noop := func(pomodoro.Interval) {}

			if err := res.Start(context.Background(), config,
				noop, noop, noop); err != nil {
				t.Fatal(err)
			}

			if res.Category != expCategory {
				t.Errorf("Expected category %q, got %q.\n",
					expCategory, res.Category)
			}

			if res.PlannedDuration != expDuration {
				t.Errorf("Expected PlannedDuration %q, got %q.\n",
					expDuration, res.PlannedDuration)
			}

			if res.State != pomodoro.StateNotStarted {
				t.Errorf("Expected State = %q, got %q.\n",
					pomodoro.StateNotStarted, res.State)
			}

			ui, err := repo.ByID(res.ID)
			if err != nil {
				t.Errorf("Expected no error. Got %q.\n", err)
			}

			if ui.State != pomodoro.StateDone {
				t.Errorf("Expected State = %q, got %q.\n",
					pomodoro.StateDone, res.State)
			}
		})
	}
}

func TestPause(t *testing.T) {
	const duration = 2 * time.Second

	repo, cleanup := getRepo(t)
	defer cleanup()

	config := pomodoro.NewConfig(repo, duration, duration, duration)

	testCases := []struct {
		name        string
		start       bool
		expState    int
		expDuration time.Duration
	}{
		{
			name: "NotStarted", start: false,
			expState: pomodoro.StateNotStarted,
			expDuration: 0,
		},
		{
			name: "Paused",
			start: true,
			expState: pomodoro.StatePaused,
			expDuration: duration/2,
		},
	}
	expError := pomodoro.ErrIntervalNotRunning

	// Execute tests for Pause
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			i, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}

			start := func(pomodoro.Interval) {}
			end := func(pomodoro.Interval) {
				t.Errorf("End callback should not be executed")
			}
			periodic := func(i pomodoro.Interval) {
				if err := i.Pause(config); err != nil {
					t.Fatal(err)
				}
			}

			if tc.start {
				if err := i.Start(ctx, config, start, periodic, end); err != nil {
					t.Fatal(err)
				}
			}

			i, err = pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}

			err = i.Pause(config)
			if err != nil {
				if !errors.Is(err, expError) {
					t.Fatalf("Expected error %q, got %q", expError, err)
				}
			}

			if err == nil {
				t.Errorf("Expected error %q, got nil", expError)
			}

			i, err = repo.ByID(i.ID)
			if err != nil {
				t.Fatal(err)
			}

			if i.State != tc.expState {
				t.Errorf("Expected state %d, got %d.\n", 
					tc.expState, i.State)

			}

			if i.ActualDuration != tc.expDuration {
				t.Errorf("Expected duration %q, got %q.\n", 
					tc.expDuration, i.ActualDuration)	
			}
			cancel()
		})
	}
}

func TestStart(t *testing.T) {
	const duration = 2 * time.Second

	repo, cleanup := getRepo(t)
	defer cleanup()

	config := pomodoro.NewConfig(repo, duration, duration, duration)

	testCases := []struct {
		name		string
		cancel 		bool
		expState	int
		expDuration time.Duration
	}{
		{
			name: "Finish", cancel: false,
			expState: pomodoro.StateDone, expDuration: duration,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())

			i, err := pomodoro.GetInterval(config)
			if err != nil {
				t.Fatal(err)
			}

			start := func(i pomodoro.Interval) {
				if i.State != pomodoro.StateRunning {
					t.Errorf("Expected state %d, got %d.\n", 
						pomodoro.StateRunning, i.State)
				}
				if i.ActualDuration >= i.PlannedDuration {
					t.Errorf("Expected ActualDuration %q, less than Planned %q.\n", 
						i.ActualDuration, i.PlannedDuration)
				}
			}

			end := func(i pomodoro.Interval) {
				if i.State != tc.expState {
					t.Errorf("Expected state %d, got %d.\n", 
						tc.expState, i.State)
				}
				if tc.cancel {
					t.Errorf("Enc callback should not be executed")
				}
			}

			periodic := func(i pomodoro.Interval) {
				if i.State != pomodoro.StateRunning {
					t.Errorf("Expected state %d, got %d.\n",
					pomodoro.StateRunning, i.State)
				}
				if tc.cancel {
					cancel()
				}
			}

			if err := i.Start(ctx, config, start, periodic, end); err != nil {
				t.Fatal(err)
			}

			i, err = repo.ByID(i.ID)
			if err != nil {
				t.Fatal(err)
			}

			if i.State != tc.expState {
				t.Errorf("Expected state %d, got %d.\n", 
					tc.expState, i.State)
			}

			if i.ActualDuration != tc.expDuration {
				t.Errorf("Expected ActualDuration %q, got %q.\n", 
					tc.expDuration, i.ActualDuration)	
			}
			cancel()
		})
	}
}
