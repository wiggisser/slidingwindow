//Package slidingwindow provides a simple implementation of requestlimits in a sliding time window.
package slidingwindow

import (
	"fmt"
	"math"
	"sync"
	"time"
)

//Limit a sliding window request limitation counter
type Limit interface {
	//Check checks whether the needed number of usages is still available in the current window
	Check(needed int) bool
	//Reset resets the currentallowance to full quota
	Reset()
}

type limit struct {
	//amount of allowed usage per timeframe
	allowance float64
	//size of the timeframe in seconds
	windowsize float64

	//currently still allowed usage in the sliding window
	currentallowance float64
	//when was the limit last checked
	lastcheck time.Time

	//a mutex for synchronizing the access to this limit
	mutex *sync.Mutex
}

func (l *limit) Reset() {
	if l.allowance <= 0 {
		return
	}

	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.lastcheck = time.Now().UTC()
	l.currentallowance = l.allowance
}

func (l *limit) Check(needed int) bool {
	//allowance of 0 means unlimited access
	if l.allowance <= 0 {
		return true
	}

	//need to synchronize the check
	l.mutex.Lock()
	defer l.mutex.Unlock()

	fneeded := float64(needed)

	//if one needs more than allowed in the whole window, request can't be fulfilled
	if fneeded > l.allowance {
		return false
	}

	now := time.Now().UTC()
	diff := now.Sub(l.lastcheck)
	l.lastcheck = now
	//update the current allowance with respect to the time that has passed since the last check
	l.currentallowance = math.Min(l.allowance, l.currentallowance+diff.Seconds()*l.allowance/l.windowsize)

	//if not enough allowance is available, reject
	if l.currentallowance < fneeded {
		return false
	}

	//decrease the current allowance by the consumed requests and approve
	l.currentallowance -= fneeded
	return true
}

var (
	namedLimits = make(map[string]Limit)
	mutex       = &sync.Mutex{}
)

//NewLimit creates a new limit with given allowance and windowsize. Use allowance = 0 for unlimited access
func NewLimit(allowance int, windowsize int) (Limit, error) {
	if allowance > 0 && windowsize <= 0 {
		return nil, fmt.Errorf("invalid value '%d' for windowsize", windowsize)
	}

	var limit = limit{
		allowance:        float64(allowance),
		windowsize:       float64(windowsize),
		currentallowance: float64(allowance),
		lastcheck:        time.Now().UTC(),
	}

	if allowance > 0 && windowsize > 0 {
		limit.mutex = &sync.Mutex{}
	}

	return &limit, nil
}

//NewNamedLimit create a new named limit with given allowance and windowsize. Use allowance = 0 for unlimited access
func NewNamedLimit(name string, allowance int, windowsize int) error {
	if name == "" {
		return fmt.Errorf("invalid value '%s' for name", name)
	}

	mutex.Lock()
	defer mutex.Unlock()
	if _, exists := namedLimits[name]; exists {
		return fmt.Errorf("named limit '%s' already exists", name)
	}

	l, e := NewLimit(allowance, windowsize)
	if e != nil {
		return e
	}
	namedLimits[name] = l
	return nil
}

//Reset resets the current allowance of the named limit to full quota
func Reset(name string) error {
	mutex.Lock()
	defer mutex.Unlock()

	l, exists := namedLimits[name]
	if !exists {
		return fmt.Errorf("named limit '%s' does not exist", name)
	}

	l.Reset()
	return nil
}

/*
Check checks whether the needed number of usages is still available in the current window of the named limit

returns
*/
func Check(name string, needed int) (bool, error) {
	mutex.Lock()
	defer mutex.Unlock()

	l, exists := namedLimits[name]
	if !exists {
		return false, fmt.Errorf("named limit '%s' does not exist", name)
	}

	return l.Check(needed), nil
}
