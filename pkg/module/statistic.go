package module

import "sync"

type Statistic struct {
	Name string
	Data sync.Map
}

func NewStatistic(name string) *Statistic {
	return &Statistic{
		Name: name,
		Data: sync.Map{},
	}
}

func (s *Statistic) AddInt64(k interface{}, v int64) {
	o, ok := s.Data.Load(k)
	if !ok {
		s.Data.Store(k, v)
	} else {
		s.Data.Store(k, v+o.(int64))
	}
}
