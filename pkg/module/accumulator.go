package module

import "sync/atomic"

type Accumulator struct {
	Result chan int64
	Count  int64
}

func NewAccumuator() *Accumulator {
	accum := &Accumulator{
		Result: make(chan int64),
	}
	accum.run()
	return accum
}

func (a *Accumulator) run() {
	go func() {
		for {
			select {
			case c := <-a.Result:
				atomic.AddInt64(&a.Count, c)
			}
		}
	}()
}
