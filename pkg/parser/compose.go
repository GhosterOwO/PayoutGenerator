package parser

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
	"gopkg.in/yaml.v2"
)

type Service struct {
	Command     string `yaml:"command"`
	draws       int64
	replacement bool
	selections  int64
	min, max    int64
}

func (s *Service) parse() error {
	var err error
	arguments := strings.Split(s.Command, " ")
	for i := 0; i < len(arguments); i += 2 {
		key := arguments[i]
		value := arguments[i+1]
		switch key {
		case "--draws":
			if err == nil {
				s.draws, err = strconv.ParseInt(value, 10, 64)
			}
		case "--range":
			subs := strings.Split(value, ",")
			s.min, err = strconv.ParseInt(subs[0], 10, 64)
			s.max, err = strconv.ParseInt(subs[1], 10, 64)
		case "--replacement":
			s.replacement = value == "yes"
		case "--selections":
			s.selections, err = strconv.ParseInt(value, 10, 64)
		}
	}
	return err
}

func (s *Service) String() string {
	return fmt.Sprintf("Draws %v, Replacement %v, RangeMin %v, RangeMax %v, Selections %v",
		humanize.Comma(s.draws), s.replacement, s.min, s.max, humanize.Comma(s.selections),
	)
}

func (s *Service) Draws() int64 {
	return s.draws
}

func (s *Service) Replacement() bool {
	return s.replacement
}

func (s *Service) Selections() int64 {
	return s.selections
}

func (s *Service) RangeMax() int64 {
	return s.max
}

func (s *Service) RangeMin() int64 {
	return s.min
}

type Compose struct {
	Services map[string]*Service `yaml:"services"`
}

func (c *Compose) Parse(b []byte) error {
	if err := yaml.Unmarshal(b, c); err != nil {
		return err
	}
	for _, s := range c.Services {
		if err := s.parse(); err != nil {
			return err
		}
	}
	return nil
}
