Build a small URL-shortening service with POST /shorten generating a Base62 code and GET /{code} returning a 302 redirect. Start with an in-memory store, then move the data to Redis or Postgres and add hit counters.

Go is a good fit here because: the in-memory store introduces concurrent access. Multiple HTTP requests can read and write the same map at the same time, so you’ll need to protect it with sync.RWMutex. Run go test -race to catch unsafe access. This is also one of those golang projects for beginners where you can understand Go’s approach to sharing state between goroutines.

A simple mutex-guarded store could look like this; you can use this for your :

type Store struct {

    mu   sync.RWMutex

    data map[string]string

}

func (s \*Store) Get(key string) (string, bool) {

    s.mu.RLock()

    defer s.mu.RUnlock()

    value, ok := s.data[key]

    return value, ok

}

Here, RLock() allows multiple reads at the same time while preventing a write from happening concurrently. You can then use the same pattern for methods that modify the map with Lock().

Difficulty: Beginner – Intermediate

Key packages: net/http, go-chi/chi or gin-gonic/gin, sync, encoding/json, redis/go-redis

What it demonstrates: HTTP handlers and middleware, appropriate status codes, concurrent-safe shared state, and the difference between a regular Go map and one that can safely be accessed by multiple goroutines.

For the HTTP layer, you can use either Gin or Chi. Both are worth trying, but the standard library’s net/http is enough to build the core API too.
