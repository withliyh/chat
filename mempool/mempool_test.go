package mempool

import (
	"fmt"
	"testing"
)

func run(pool *MemPool) {
	for i := 0; i < 100; i++ {
		buf := pool.Get()
		pool.Give(buf)
	}
}

func Test_MemPool(t *testing.T) {
	conf := &Config{2048, 1024}
	pool := NewMemPool(conf)
	/*
		for i := 0; i < 10; i++ {
			go run(pool)
		}
	*/
	var buf []byte
	for i := 0; i < 100; i++ {
		buf = pool.Get()
	}

	for i := 0; i < 100; i++ {
		pool.Give(buf)
	}
	fmt.Println(pool.Len())

}
