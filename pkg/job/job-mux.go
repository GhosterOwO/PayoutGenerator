package job

import "sync"

var muxs map[string]*sync.Mutex

func init() {
	muxs = make(map[string]*sync.Mutex)
}

func aquireMutex(filename string) *sync.Mutex {
	mu, ok := muxs[filename]
	if !ok {
		mu = &sync.Mutex{}
		muxs[filename] = mu
	}
	return mu
}
