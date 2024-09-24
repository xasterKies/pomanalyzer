// +build inmemory

package cmd

import (
	"github.com/xasterKies/pomanalyzer/pomodoro"
	"github.com/xasterKies/pomanalyzer/repository"
)

func getRepo() (pomodoro.Repository, error) {
  return repository.NewInMemoryRepo(), nil
}