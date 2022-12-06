package job

import (
	"bytes"
	"io"
	"log"
	"strconv"
	"strings"

	"future.net.co/rngcollector/pkg/module"
)

// Remove line which container element 0.
type OccurrenceJob struct {
	Statistic *module.Statistic
	Chunk     *bytes.Buffer
}

func NewOccurrenceJob(s *module.Statistic, b *bytes.Buffer) Job {
	c := bufferPool.Get()
	chunk := c.(*bytes.Buffer)
	chunk.ReadFrom(b)
	return &OccurrenceJob{
		Statistic: s,
		Chunk:     chunk,
	}
}

func (j *OccurrenceJob) Execute() {
	c := bufferPool.Get()
	chunk := c.(*bytes.Buffer)
	defer chunk.Reset()
	defer bufferPool.Put(chunk)
	defer j.Chunk.Reset()
	defer bufferPool.Put(j.Chunk)
	for {
		line, err := j.Chunk.ReadString('\n')
		s := strings.Trim(line, "\n\r")
		a := strings.Split(strings.TrimSpace(s), " ")
		for _, v := range a {
			if i, err := strconv.ParseUint(v, 10, 64); err == nil {
				j.Statistic.AddInt64(i, 1)
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Print(err)
			break
		}
	}
}
