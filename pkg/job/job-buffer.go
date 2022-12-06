package job

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{New: func() interface{} {
	buffer := &bytes.Buffer{}
	return buffer
}}
