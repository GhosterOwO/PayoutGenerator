package random

import (
	"bufio"
	"encoding/binary"
	"log"
	"math/rand"
	"os"
	"sync"
	"sync/atomic"

	"github.com/seehuhn/mt19937"
)

type Generator struct {
	lk *sync.Mutex
	g  *rand.Rand
	c  int32
}

func NewGenerator() *Generator {
	g := rand.New(mt19937.New())
	g.Seed(urandom())
	return &Generator{
		lk: &sync.Mutex{},
		g:  g,
	}
}

func urandom() int64 {
	b := make([]byte, 8)
	if f, err := os.Open("/dev/urandom"); err == nil {
		defer f.Close()
		r := bufio.NewReader(f)
		if n, err := r.Read(b); err != nil {
			log.Panic(n, err)
		}
	} else {
		if n, err := rand.Read(b); err != nil {
			log.Panic(n, err)
		}
	}
	return int64(binary.BigEndian.Uint64(b))
}

func (g *Generator) Intn(n int) int {
	g.lk.Lock()
	defer g.lk.Unlock()
	c := atomic.LoadInt32(&g.c)
	if c > 0 && c%624 == 0 {
		g.g.Seed(urandom())
	}
	atomic.AddInt32(&g.c, 1)
	return g.g.Intn(n)
}

func (g *Generator) Shuffle(a []int) {
	g.lk.Lock()
	defer g.lk.Unlock()
	swap := func(i, j int) {
		a[i], a[j] = a[j], a[i]
	}
	g.g.Shuffle(len(a), swap)
}
