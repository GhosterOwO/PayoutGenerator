package util

import (
	"bufio"
	"bytes"
	"io"
	"log"
	"os"
	"sync"
)

var chunkPool = sync.Pool{New: func() interface{} {
	buffer := &bytes.Buffer{}
	buffer.Grow(MaxReaderChunk)
	return buffer
}}

var MaxReaderChunk = 8 * 1024

type reader func(*bytes.Buffer)

func ReadChunk(filename string, handle reader) {
	f, err := os.Open(filename)
	if err != nil {
		log.Print(err)
		return
	}
	defer f.Close()
	c := chunkPool.Get()
	defer chunkPool.Put(c)
	chunk := c.(*bytes.Buffer)
	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		chunk.WriteString(line)
		if chunk.Len() >= MaxReaderChunk {
			if chunk.Len() > 0 {
				handle(chunk)
				chunk.Reset()
			}
		}
		if err == io.EOF {
			break
		} else if err != nil {
			log.Print(err)
			break
		}
	}
	if chunk.Len() > 0 {
		handle(chunk)
		chunk.Reset()
	}
}
