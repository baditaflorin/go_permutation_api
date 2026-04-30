package api

import (
	"fmt"
	"net/http"
	"runtime/debug"
)

// recoveryMiddleware recovers from panics and returns 500 error
func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic
				fmt.Printf("PANIC: %v\n", err)
				fmt.Printf("Stack trace:\n%s\n", debug.Stack())

				// Return 500 error
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)

				// Increment error counter
				IncrementErrors()
			}
		}()

		next.ServeHTTP(w, r)
	})
}
