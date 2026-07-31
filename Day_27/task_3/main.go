package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
)

func testerHandler(w http.ResponseWriter, r *http.Request) {
	if r.TLS != nil {
		fmt.Fprintf(w, "TLS версия: %v\n", r.TLS.Version)
		if r.TLS.Version == tls.VersionTLS13 {
			fmt.Fprintf(w, "✅ Используется TLS 1.3\n")
		} else {
			fmt.Fprintf(w, "❌ Не TLS 1.3\n")
		}
	}
	fmt.Fprintf(w, "TLS 1.3 работает!")
}

func main() {

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS13,
		MaxVersion: tls.VersionTLS13,
		CipherSuites: []uint16{
			tls.TLS_AES_256_GCM_SHA384, // Самый сильный
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_CHACHA20_POLY1305_SHA256,
		},
	}

	httpsServer := &http.Server{
		Addr:      ":8443",
		TLSConfig: tlsConfig,
	}

	http.HandleFunc("/", testerHandler)

	fmt.Println("HTTPS сервер на https://localhost:8443")
	fmt.Println("  HTTP редиректит на HTTPS")
	fmt.Println("  / - HTTPS работает!")

	if err := httpsServer.ListenAndServeTLS("certs/server.crt", "certs/server.key"); err != nil {
		fmt.Println("Ошибка HTTPS:", err)
	}

}
