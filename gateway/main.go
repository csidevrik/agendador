package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

var secret = []byte("secret")

func createToken(username string) (string, error) {
	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}
	headerJSON, _ := json.Marshal(header)

	claims := map[string]interface{}{
		"sub": username,
		"exp": time.Now().Add(time.Hour).Unix(),
	}
	claimsJSON, _ := json.Marshal(claims)

	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)
	unsigned := headerEncoded + "." + claimsEncoded

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(unsigned))
	signature := h.Sum(nil)
	signatureEncoded := base64.RawURLEncoding.EncodeToString(signature)

	token := unsigned + "." + signatureEncoded
	return token, nil
}

func verifyToken(token string) (string, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", false
	}
	unsigned := parts[0] + "." + parts[1]

	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return "", false
	}

	h := hmac.New(sha256.New, secret)
	h.Write([]byte(unsigned))
	expected := h.Sum(nil)

	if !hmac.Equal(sig, expected) {
		return "", false
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", false
	}
	var claims struct {
		Sub string `json:"sub"`
		Exp int64  `json:"exp"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return "", false
	}
	if time.Now().Unix() > claims.Exp {
		return "", false
	}

	return claims.Sub, true
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	json.Unmarshal(body, &creds)

	if creds.Username == "admin" && creds.Password == "password" {
		token, _ := createToken(creds.Username)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"token": token})
		return
	}

	http.Error(w, "unauthorized", http.StatusUnauthorized)
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		user, ok := verifyToken(token)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		r.Header.Set("X-User", user)
		next.ServeHTTP(w, r)
	})
}

func proxyHandler(target string) http.Handler {
	url, err := url.Parse(target)
	if err != nil {
		log.Fatalf("invalid proxy target: %v", err)
	}
	proxy := httputil.NewSingleHostReverseProxy(url)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/login", loginHandler)

	apiProxy := proxyHandler("http://localhost:8000")
	mux.Handle("/api/", jwtMiddleware(apiProxy))

	log.Println("Gateway listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
