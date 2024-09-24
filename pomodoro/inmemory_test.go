// +build inmemory

package pomodoro_test

import (
	"testing"

	"github.com/xasterKies/pomanalyzer/pomodoro"
	"github.com/xasterKies/pomanalyzer/repository"
)

func getRepo(t *testing.T) (pomodoro.Repository, func()) {
	t.Helper()

	return repository.NewInMemoryRepo(), func() {}
}