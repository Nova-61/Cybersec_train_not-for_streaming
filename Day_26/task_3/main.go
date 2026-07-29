package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "HTTPS работает!")
}

func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://localhost:8443"+r.URL.Path, http.StatusMovedPermanently)
}

func main() {
	// Настройка TLS
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		},
	}

	// HTTP сервер (редирект на HTTPS)
	go func() {
		fmt.Println("HTTP сервер на порту 8080 → редирект на HTTPS")
		if err := http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS)); err != nil {
			fmt.Println("Ошибка HTTP:", err)
		}
	}()

	// HTTPS сервер
	httpsServer := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", helloHandler)

	fmt.Println("HTTPS сервер на https://localhost:8443")
	fmt.Println("  HTTP редиректит на HTTPS")
	fmt.Println("  / - HTTPS работает!")

	if err := httpsServer.ListenAndServeTLS("certs/server.crt", "certs/server.key"); err != nil {
		fmt.Println("Ошибка HTTPS:", err)
	}
}
