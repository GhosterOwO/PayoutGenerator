package parser

import (
	"testing"

	"future.net.co/rngcollector/assets"
)

func TestDecodeCompose(t *testing.T) {
	c := &Compose{}
	if b, err := assets.Asset("assets/docker-compose.yml"); err == nil {
		if err := c.Parse(b); err != nil {
			t.Error(err)
		}
	}
	for _, s := range c.Services {
		if s.RangeMin() == 0 {
			t.Errorf("Error %v", s)
		} else if s.RangeMax() == 0 {
			t.Errorf("Error %v", s)
		} else if s.Selections() == 0 {
			t.Errorf("Error %v", s)
		}
	}
}
