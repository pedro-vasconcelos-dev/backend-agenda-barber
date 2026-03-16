package cookies

import (
	"fmt"
	"net/http"
	"time"
)

// Ajuste conforme seu domínio e tempo de vida
type CookieCfg struct {
	Domain        string        // ex: ".legitimatech.com.br"
	RefreshName   string        // "refresh_token"
	CSRFName      string        // "csrf_token"
	RefreshMaxAge time.Duration // ex: 30 * 24h
	CSRFSameSite  http.SameSite // normalmente SameSite=Lax para csrf cookie
}

// Set refresh cookie (HttpOnly, Secure, SameSite=None)
func SetRefreshCookie(w http.ResponseWriter, cfg CookieCfg, value string) {
	maxAge := int(cfg.RefreshMaxAge.Seconds())
	// Como SameSite=None não está em http.Cookie nas versões antigas, escrevemos manualmente:
	cookie := fmt.Sprintf("%s=%s; Path=/; Max-Age=%d; Domain=%s; HttpOnly; Secure; SameSite=None",
		cfg.RefreshName, value, maxAge, cfg.Domain)
	w.Header().Add("Set-Cookie", cookie)
}

// Apaga refresh cookie
func ClearRefreshCookie(w http.ResponseWriter, cfg CookieCfg) {
	cookie := fmt.Sprintf("%s=; Path=/; Max-Age=0; Domain=%s; HttpOnly; Secure; SameSite=None",
		cfg.RefreshName, cfg.Domain)
	w.Header().Add("Set-Cookie", cookie)
}

// CSRF cookie (NÃO HttpOnly, para o front ler e mandar no header X-CSRF-Token)
func SetCSRFCookie(w http.ResponseWriter, cfg CookieCfg, value string, ttl time.Duration) {
	http.SetCookie(w, &http.Cookie{
		Name:     cfg.CSRFName,
		Value:    value,
		Path:     "/",
		Domain:   cfg.Domain,
		MaxAge:   int(ttl.Seconds()),
		SameSite: cfg.CSRFSameSite,
		Secure:   true,
		HttpOnly: false, // importante: false
	})
}
