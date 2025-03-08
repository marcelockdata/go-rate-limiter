package limiter

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

type RateLimiter struct {
	store          Store
	rateLimitIP    int
	rateLimitToken int
	blockDuration  time.Duration
}

func NewRateLimiter(store Store) *RateLimiter {
	rateLimitIP, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_IP"))
	rateLimitToken, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_TOKEN"))
	blockDuration, _ := strconv.Atoi(os.Getenv("BLOCK_DURATION"))

	return &RateLimiter{
		store:          store,
		rateLimitIP:    rateLimitIP,
		rateLimitToken: rateLimitToken,
		blockDuration:  time.Duration(blockDuration) * time.Second,
	}
}

func getIP(r *http.Request) string {
	headers := []string{"X-Forwarded-For", "X-Real-IP"}
	for _, header := range headers {
		ip := r.Header.Get(header)
		if ip != "" {
			return strings.Split(ip, ",")[0]
		}
	}
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

func (rl *RateLimiter) CheckRateLimit(ip, token string) (bool, error) {
	var limit int
	if token != "" {
		limit = rl.rateLimitToken
	} else {
		limit = rl.rateLimitIP
	}

	key := ip
	if token != "" {
		key = token
	}

	return rl.store.Allow(key, limit, rl.blockDuration)
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		token := r.Header.Get("API_KEY")

		allowed, err := rl.CheckRateLimit(ip, token)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if !allowed {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func InitializeRateLimiters() *RateLimiter {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
		return nil
	}

	redis_store := NewRedisStore(rdb)
	return NewRateLimiter(redis_store)
}
