package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"sync"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {

	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client)
	)

	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()

		}
	}()

	// Create a new rate limiter limiter
	// limiter := rate.NewLimiter(2, 4)
	// Create a new handler that wraps the next handler
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {
			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				// app.rateLimitExceededResponse(w, r, err)
				app.serverErrorResponse(w, r, err)
				return
			}
			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.config.limiter.rps), app.config.limiter.burst),
				}
			}

			// ip, _, err := net.SplitHostPort(r.RemoteAddr)
			// if err != nil {
			// 	// app.rateLimitExceededResponse(w, r, err)
			// 	app.serverErrorResponse(w, r, err)
			// 	return
			// }
			// mu.Lock()

			// if _, found := clients[ip]; !found {
			// 	clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
			// }

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}

		next.ServeHTTP(w, r)
	})
}
