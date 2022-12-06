package job

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"
)

var MaxRotateSize = int64(64)
var writers map[string]*RotateWriter
var initial = time.Now().Format("20060102150405")

func init() {
	writers = make(map[string]*RotateWriter)
}

func NewRotateWriter(filename string) *RotateWriter {
	w, ok := writers[filename]
	if !ok {
		w = (&RotateWriter{filename: filename, folder: initial})
		writers[filename] = w
	}
	return w
}

type RotateWriter struct {
	lock     sync.Mutex
	counts   uint64
	folder   string
	filename string
}

// Write satisfies the io.Writer interface.
func (w *RotateWriter) Write(b []byte) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.rotate()
	const flags = os.O_APPEND | os.O_WRONLY | os.O_CREATE
	os.MkdirAll(filepath.Dir(w.filename), os.ModePerm)
	f, err := os.OpenFile(w.filename, flags, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return err
}

func makeFilename(f string, i int) string {
	path := filepath.Join(filepath.Dir(f), strconv.Itoa(i))
	new := filepath.Join(path, filepath.Base(f))
	if _, err := os.Stat(new); errors.Is(err, os.ErrNotExist) {
		return new
	}
	return makeFilename(f, i+1)
}

func (w *RotateWriter) rotate() {
	fi, err := os.Stat(w.filename)
	if err != nil {
		return
	}
	if fi.Size() <= 1024*1024*MaxRotateSize {
		return
	}
	filename := makeFilename(w.filename, 0)
	os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err := os.Rename(w.filename, filename); err != nil {
		log.Print(err)
	}
}

type WriteJob struct {
	Result string
	writer *RotateWriter
}

func NewWriteJob(filename string, result string) Job {
	return &WriteJob{
		writer: NewRotateWriter(filename),
		Result: result,
	}
}

func (j *WriteJob) Execute() {
	if err := j.writer.Write([]byte(j.Result)); err != nil {
		log.Print(err)
	}
}
