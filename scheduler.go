package main

type scheduler interface {
	Run([]func()) chan bool
}

// Scheduler ...
type Scheduler struct {
}

// Run ...
func (s *Scheduler) Run(fns []func()) {
	chn := make(chan bool, len(fns))
	for _, fn := range fns {
		go func(fn func()) {
			fn()
			chn <- true
		}(fn)
	}
	for i := 0; i < len(fns); i++ {
		<-chn
	}
}
