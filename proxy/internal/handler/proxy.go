package handler

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ProxyHandler struct {
	proxy  *httputil.ReverseProxy
	target *url.URL
}

func NewProxyHandler(targetURL string) *ProxyHandler {
	// Парсим целевой URL
	target, err := url.Parse(targetURL)
	if err != nil {
		panic("Invalid target URL: " + err.Error())
	}

	// Создаём ReverseProxy
	proxy := httputil.NewSingleHostReverseProxy(target)

	// Кастомизируем запрос
	proxy.Director = func(req *http.Request) {
		// Меняем хост на целевой
		req.Host = target.Host
		req.URL.Host = target.Host
		req.URL.Scheme = target.Scheme

		// Добавляем заголовки
		req.Header.Add("X-Proxy", "true")
		req.Header.Add("X-Forwarded-For", req.RemoteAddr)
	}

	// Обработка ошибок
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		fmt.Printf("[ERROR] %v\n", err)
		http.Error(w, "Proxy error: "+err.Error(), http.StatusBadGateway)
	}

	return &ProxyHandler{
		proxy:  proxy,
		target: target,
	}
}

func (p *ProxyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	// Логируем запрос
	fmt.Printf("[PROXY] %s %s -> %s%s\n",
		r.Method,
		r.URL.Path,
		p.target.Host,
		r.URL.Path)

	// Передаём запрос прокси
	p.proxy.ServeHTTP(w, r)

	// Логируем время
	fmt.Printf("[PROXY] Completed in %v\n", time.Since(start))
}
