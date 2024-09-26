package pomodoro

import (
  "fmt"
  "time"
)

type LineSeries struct {
  Name   string
  Labels map[int]string
  Values []float64
}

func DailySummary(day time.Time,
  config *IntervalConfig) ([]time.Duration, error) {

  dPomo, err := config.repo.CategorySummary(day, CategoryPomodoro)
  if err != nil {
    return nil, err
  }

  dBreaks, err := config.repo.CategorySummary(day, "%Break")
  if err != nil {
    return nil, err
  }

  return []time.Duration{
    dPomo,
    dBreaks,
  }, nil
}

func RangeSummary(start time.Time, n int,
  config *IntervalConfig) ([]LineSeries, error) {

  pomodoroSeries := LineSeries{
    Name:   "Pomodoro",
    Labels: make(map[int]string),
    Values: make([]float64, n),
  }

  breakSeries := LineSeries{
    Name:   "Break",
    Labels: make(map[int]string),
    Values: make([]float64, n),
  }

  for i := 0; i < n; i++ {
    day := start.AddDate(0, 0, -i)
    ds, err := DailySummary(day, config)
    if err != nil {
      return nil, err
    }

    label := fmt.Sprintf("%02d/%s", day.Day(), day.Format("Jan"))

    pomodoroSeries.Labels[i] = label
    pomodoroSeries.Values[i] = ds[0].Seconds()

    breakSeries.Labels[i] = label
    breakSeries.Values[i] = ds[1].Seconds()
  }

  return []LineSeries{
    pomodoroSeries,
    breakSeries,
  }, nil
}