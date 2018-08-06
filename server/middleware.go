package server

import "net/http"

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
