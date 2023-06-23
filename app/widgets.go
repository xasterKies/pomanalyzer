package app

import (
	"context"

	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/widgets/donut"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/mum/termdash/widgets/text"
)

type widgets struct {
	donTimer	*donut.Donut
	disType		*segmentdisplay.SegmentDisplay
	txtInfo		*text.Text
	txtTimer	*text.Text
	updateDonTimer	chan []int
	updateTxtInfo	chan string
	updateTxtTimer	chan string
	updateTxtType	chan string
}

func (w *widgets) update(timer []int, txtType, txtInfo, txtTimer string,
	redrawCh chan<- bool) {

		if txtInfo != "" {
			w.updateTxtInfo <- txtInfo
		}

		if txtType != "" {
			w.updateTxtType <- txtType
		}

		if txtTimer != "" {
			w.updateTxtTimer <- txtTimer
		}

		if len(timer) > 0 {
			w.updateDonTimer <- timer
		}

		redrawCh <- true
	}

func newWidgets(ctx context.Context, errorCh chan<- error) (*widgets, error) {
	w := &widgets{}
	var err error

	w.updateDonTimer = make(chan []int)
	w.updateTxtType = make(chan string)
	w.updateTxtInfo = make(chan string)
	w.updateTxtTimer =  make(chan string)

	w.donTimer, err = newDonut(ctx, w.updateDonTimer, errorCh)
	if err != nil {
		return nil, err
	}

	w.disType, err = newSegmentDisplay(ctx, w.updateTxtType, errorCh)
	if err != nil {
		return nil, err
	}

	w.txtInfo, err = newText(ctx, w.updateTxtInfo, errorCh)
	if err != nil {
		return nil, err
	}

	w.txtTimer, err = newText(ctx, w.updateTxtTimer, errorCh)
	if err != nil {
		return nil, err
	}

	return w, nil
}