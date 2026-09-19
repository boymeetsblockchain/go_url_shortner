package test

import (
	"sync"
	"testing"

	"github.com/boymeetsblockchain/url_shortner/internal/shortener"
	"github.com/boymeetsblockchain/url_shortner/internal/store"
)

func TestConcurrentAccess(t *testing.T) {
	s := store.New()
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			code := shortener.ToBase62(n)
			s.Set(code, "https://example.com")
			s.Get(code)
			s.IncrementHits(code)
		}(i)
	}
	wg.Wait()
}
