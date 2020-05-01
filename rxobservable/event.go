package rxobservable

import "sync"

type eventStream interface {
	On(string, func(interface{}) (interface{}, error))
	Trigger(string, interface{})
}

type stream struct {
	Stream map[string]chan interface{}
	M      sync.RWMutex
}

func (ev *stream) Trigger(eventType string, value interface{}) {
	ev.M.Lock()
	ev.Stream[eventType] <- value
	ev.M.Unlock()
}

func (ev *stream) On(eventType string, cb func(interface{}) (interface{}, error)) (interface{}, error) {
	for evType, value := range ev.Stream {
		if evType == eventType {
			ev.M.RLock()
			defer ev.M.RUnlock()
			return cb(<-value)
		}
	}
	return nil, nil
}
