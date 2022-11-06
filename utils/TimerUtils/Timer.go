package TimerUtils

import (
	"errors"
	"fmt"
	"time"
)

var ST = &simpleTimer{timers: make(map[string]*time.Timer)}

type simpleTimer struct {
	timers map[string]*time.Timer
}

func (s *simpleTimer) StartTimer(key string, d time.Duration, f func()) {
	if t, ok := s.timers[key]; ok {
		t.Stop()
	}
	s.timers[key] = time.NewTimer(d)
	go func() {
		select {
		case <-s.timers[key].C:
			f()
			delete(s.timers, key)
		}
	}()
}

func (s *simpleTimer) StopTimer(key string) (bool, error) {
	if _, ok := s.timers[key]; !ok {
		return false, errors.New("Not Find Timer ")
	}
	ok := s.timers[key].Stop()
	delete(s.timers, key)
	return ok, nil
}

func (s *simpleTimer) RestTimer(key string, d time.Duration) (bool, error) {
	if _, ok := s.timers[key]; !ok {
		return false, errors.New("Not Find Timer ")
	}
	return s.timers[key].Reset(d), nil
}

func (s *simpleTimer) Pln() {
	fmt.Println(s.timers)
}
