package app

import (
	"context"
	"image"
	"time"

	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/xasterKies/pomanalyzer/pomodoro"
)

type App struct {
  ctx        context.Context
  controller *termdash.Controller
  redrawCh   chan bool
  errorCh    chan error
  term       *tcell.Terminal
  size       image.Point
}

func New(config *pomodoro.IntervalConfig) (*App, error) {
  ctx, cancel := context.WithCancel(context.Background())

  quitter := func(k *terminalapi.Keyboard) {
    if k.Key == 'q' || k.Key == 'Q' {
      cancel()
    }
  }

  redrawCh := make(chan bool)
  errorCh := make(chan error)

  w, err := newWidgets(ctx, errorCh)
  if err != nil {
    return nil, err
  }

  s, err := newSummary(ctx, config, redrawCh, errorCh)
  if err != nil {
    return nil, err
  }

  b, err := newButtonSet(ctx, config, w, s, redrawCh, errorCh)
  if err != nil {
    return nil, err
  }

  term, err := tcell.New()
  if err != nil {
    return nil, err
  }

  c, err := newGrid(b, w, s, term)
  if err != nil {
    return nil, err
  }

  controller, err := termdash.NewController(term, c,
    termdash.KeyboardSubscriber(quitter))
  if err != nil {
    return nil, err
  }

  return &App{
    ctx:        ctx,
    controller: controller,
    redrawCh:   redrawCh,
    errorCh:    errorCh,
    term:       term,
  }, nil
}

func (a *App) resize() error {
  if a.size.Eq(a.term.Size()) {
    return nil
  }

  a.size = a.term.Size()
  if err := a.term.Clear(); err != nil {
    return err
  }

  return a.controller.Redraw()
}

func (a *App) Run() error {
  defer a.term.Close()
  defer a.controller.Close()

  ticker := time.NewTicker(2 * time.Second)
  defer ticker.Stop()

  for {
    select {
    case <-a.redrawCh:
      if err := a.controller.Redraw(); err != nil {
        return err
      }
    case err := <-a.errorCh:
      if err != nil {
        return err
      }
    case <-a.ctx.Done():
      return nil
    case <-ticker.C:
      if err := a.resize(); err != nil {
        return err
      }
    }
  }
}