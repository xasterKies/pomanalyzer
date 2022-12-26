package repository

import (
	"fmt"
	"sync"

	"github.com/xasterKies/pomanalyzer/pomodoro"
)

type inMemoryRepo struct {
	sync.RWMutex // preventing concurent access to the inMemory datastore
	intervals []pomodoro.Interval
}

/* 
	* Instantiate a new inMemoryRepo type with an empty slice of pomodoro.interval
*/
func NewInMemoryRepo() *inMemoryRepo {
	return &inMemoryRepo{
		intervals: []pomodoro.Interval{}
	}
}

/* 
	* Saves instance of the pomodoro.Interval value to the data store
*/
func (r *inMemoryRepo) Create(i pomodoro.Interval) (int64, error) {
	r.Lock()
	defer r.Unlock()

	i.ID = int64(len(r.intervals)) + 1

	r.intervals = append(r.intervals, i)

	return i.ID, nil
}

/* 
	* Update value of an existing entry in the data store
*/
func (r *inMemoryRepo) Update(i pomodoro.Interval) error {
	r.Lock()
	defer r.Unlock()
	if i.ID == 0 {
		return fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, i.ID)
	}

	r.intervals[i.ID] = i
	return nil
}

/* 
	* Retrieve and return an item by its id
*/
func (r *inMemoryRepo) ById(id int64) (pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	i := pomodoro.Interval{}
	if id == 0 {
		return i, fmt.Errorf("%w: %d", pomodoro.ErrInvalidID, id)
	}
	i = r.intervals[id-1]
	return i, nil
}

/* 
	* Retrieve a given number of intervals of category break
*/
func (r *inMemoryRepo) Breaks(n int) ([]pomodoro.Interval, error) {
	r.RLock()
	defer r.RUnlock()
	data := []pomodoro.Interval{}
	for k := len(r.intervals) - 1; k >= 0; k-- {
		if r.intervals[k].Category == pomodoro.CategoryPomodoro {
			continue
		}

		data = append(data, r.intervals[k])

		if len(data) == n {
			return data, nil
		}
	}
}