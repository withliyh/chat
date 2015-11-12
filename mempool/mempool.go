package mempool

import "container/list"

var Pool = NewMemPool(&Config{LimitSize: 4096, LimitCount: 1024})

type Config struct {
	LimitSize  uint32
	LimitCount uint32
}

type MemPool struct {
	config *Config
	q      *list.List
	get    chan []byte
	give   chan []byte
}

func NewMemPool(conf *Config) *MemPool {
	pool := &MemPool{
		config: conf,
		q:      list.New(),
		get:    make(chan []byte),
		give:   make(chan []byte),
	}
	go pool.run()
	return pool
}

func (this *MemPool) run() {
	for {
		if this.q.Len() == 0 {
			this.q.PushFront(make([]byte, this.config.LimitSize))
		}

		e := this.q.Front()
		select {
		case b := <-this.give:
			this.q.PushFront(b)
		case this.get <- e.Value.([]byte):
			this.q.Remove(e)
		}
	}
}

func (this *MemPool) Get() []byte {
	return <-this.get
}

func (this *MemPool) Give(slice []byte) {
	this.give <- slice
}

func (this *MemPool) Len() int {
	return this.q.Len()
}
