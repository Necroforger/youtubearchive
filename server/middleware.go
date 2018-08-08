package server

import (
	"net/http"
	"net/url"
)

// CSSMiddleware is a middleware that sets the template CSS variable
func (s *Server) CSSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var css string
		if c, err := r.Cookie("css"); err == nil {
			css = c.Value
		} else {
			css = "main.css"
		}
		css = "/static/css/" + css
		r = ctxAppendTemplateVars(r, map[string]interface{}{
			"css": css,
		})
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware requires that a page is authenticated
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("password")
		if err != nil || cookie.Value != s.Password {
			http.Redirect(w, r, "/login?redirect="+url.QueryEscape(r.URL.String()), http.StatusFound)
			return
		}
		r = ctxAppendTemplateVars(r, map[string]interface{}{
			"loggedin": true,
			"redirect": r.URL.String(),
		})
		next.ServeHTTP(w, r)
	})
}
