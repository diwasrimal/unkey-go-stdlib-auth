package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	unkeygo "github.com/unkeyed/unkey-go"
	"github.com/unkeyed/unkey-go/models/components"
)

var unkeyClient *unkeygo.Unkey

// Ensures proper environment vars are set and initalizes unkey client.
func init() {
	godotenv.Load()
	unkeyRootKey := os.Getenv("UNKEY_ROOT_KEY")
	unkeyAPIID := os.Getenv("UNKEY_API_ID")
	if unkeyRootKey == "" || unkeyAPIID == "" {
		panic("Environment variable UNKEY_ROOT_KEY and/or UNKEY_API_ID not set")
	}
	unkeyClient = unkeygo.New(unkeygo.WithSecurity(unkeyRootKey))
}

func main() {
	// Map routes to handlers
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Hello there!"))
	})

	// Apply the middleware
	handler := useUnkeyAPIAuth(mux)

	addr := ":3030"
	log.Printf("Server listening on %v...\n", addr)
	log.Fatal(http.ListenAndServe(addr, handler))
}

// Validates API key provided in request header 'Authorization'
func useUnkeyAPIAuth(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		authKey := r.Header.Get("Authorization")
		if authKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized: No API key provided"))
			return
		}

		// Request to verify api key provided in header
		req := components.V1KeysVerifyKeyRequest{
			APIID: unkeygo.String(os.Getenv("UNKEY_API_ID")),
			Key:   authKey,
		}
		res, err := unkeyClient.Keys.VerifyKey(context.Background(), req)
		if err != nil {
			log.Printf("error: %T verifying key: %v", err, err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error"))
			return
		}
		if res.V1KeysVerifyKeyResponse != nil && res.V1KeysVerifyKeyResponse.Valid {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Unauthorized: Invalid API key"))
		}
	}
	return http.HandlerFunc(fn)
}
