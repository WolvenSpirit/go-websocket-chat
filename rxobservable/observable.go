package rxobservable

import "sync"

/*
rxobservable is meant to be a proof of concept demonstrating that the observer/subscriber pattern
does not have to be different or complicated especially in a language such as Golang.
By RX standards this would be an implementation of a hot observable that broadcasts the events to all subscribers.

If Subscribe() would return an object it would need to be type asserted later,
by passing a channel multiple subscribers can easily get events sent only once
while letting you declare the type of data that you need to subscribe for.
*/
type observable interface {
	Next(interface{})
	Subscribe(func(*chan interface{}))
}

type rxobservable struct {
	buffer      chan interface{}
	subscribers map[int]*chan interface{}
	m           sync.Mutex
}

// Observable represents a stream to which you can subscribe with a channel for events (values)
type Observable struct {
	rxobservable
}

// Next - pass a new value to be broadcasted.
func (o *Observable) Next(sl ...interface{}) {
	o.m.Lock()
	for _, v := range sl {
		for _, ch := range o.subscribers {
			*ch <- v

		}
	}
	o.m.Unlock()
}
func (o *Observable) allocate(id int) {
	// ...
}

// Subscribe with your own channel to the events of the observable.
func (o *Observable) Subscribe(ch *chan interface{}) {
	o.m.Lock()
	if o.subscribers == nil {
		o.subscribers = make(map[int]*chan interface{})
	}
	o.subscribers[len(o.subscribers)+1] = ch
	o.allocate(len(o.subscribers) + 1)
	o.m.Unlock()
}

// Close - subscribed channels will be closed by using this call.
func (o *Observable) Close() {
	for _, ch := range o.subscribers {
		close(*ch)
	}
}
