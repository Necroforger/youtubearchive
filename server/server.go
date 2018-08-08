package server

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/Necroforger/youtubearchive"
	"github.com/go-chi/chi"

	"github.com/jinzhu/gorm"
)

// Server serves the database
type Server struct {
	DB       *gorm.DB
	router   *chi.Mux
	Logger   io.Writer
	Password string

	tplmu     sync.RWMutex
	templates *template.Template
}

// Log logs text to the Logger
func (s *Server) Log(data ...interface{}) {
	if s.Logger == nil {
		return
	}
	fmt.Fprintln(s.Logger, data...)
}

// Options are optional server parameters
type Options struct {
	Password string
}

// NewOptions ...
func NewOptions() *Options {
	return &Options{
		Password: "",
	}
}

// NewServer creates and initializes a new server
func NewServer(DB *gorm.DB, opts *Options) *Server {
	if opts == nil {
		opts = NewOptions()
	}
	youtubearchive.InitDB(DB)
	s := &Server{
		DB:       DB,
		router:   chi.NewMux(),
		Password: opts.Password,
	}
	s.route(s.router)
	return s
}

// LoadTemplates loads the templates into the server
func (s *Server) LoadTemplates(data ...string) error {
	s.tplmu.Lock()
	defer s.tplmu.Unlock()

	tpl := template.New("").Funcs(s.templateFuncs())
	for _, v := range data {
		_, err := tpl.Parse(v)
		if err != nil {
			return err
		}
	}
	s.templates = tpl
	return nil
}

// LoadTemplatesGlob loads templates from the file system
func (s *Server) LoadTemplatesGlob(globs ...string) error {
	s.tplmu.Lock()
	defer s.tplmu.Unlock()

	tpl := template.New("").Funcs(s.templateFuncs())
	for _, v := range globs {
		_, err := tpl.ParseGlob(v)
		if err != nil {
			return err
		}
	}
	s.templates = tpl
	return nil
}

// templateFuncs returns default template functions
func (s *Server) templateFuncs() template.FuncMap {
	return template.FuncMap{
		"escape_spaces": func(text string) string {
			return strings.Replace(text, " ", `\ `, -1)
		},
	}
}

// route creates routes
func (s *Server) route(r *chi.Mux) {
	r.Use(s.CSSMiddleware) // Set the custom css

	r.Get("/login", s.LoginHandler)
	r.Post("/login", s.LoginHandlerPost)

	r.Group(func(r chi.Router) {
		if s.Password != "" {
			r.Use(s.AuthMiddleware)
			r.Get("/logout", s.LogoutHandler)
		}
		r.Get("/search", s.HandleSearch)
		r.Get("/view", s.HandleView)
		r.Get("/channels", s.HandleChannels)
		r.Get("/", s.HandleHome)
		r.HandleFunc("/*", s.HandleHome)
	})
}

// GetRoutes returns this router
func (s *Server) GetRoutes() http.Handler {
	return s.router
}

// ExecuteTemplate executes a template
func (s *Server) ExecuteTemplate(w http.ResponseWriter, r *http.Request, name string, data map[string]interface{}) error {
	s.tplmu.RLock()
	defer s.tplmu.RUnlock()

	if extra := ctxGetTemplateVars(r); extra != nil {
		extendTemplateVars(*extra, &data)
	}

	err := s.templates.ExecuteTemplate(w, name, data)
	if err != nil {
		s.Log("error executing template ", name, ": ", err.Error())
	}
	return err
}

func extendTemplateVars(a map[string]interface{}, b *map[string]interface{}) {
	for key := range a {
		if _, ok := (*b)[key]; !ok {
			(*b)[key] = a[key]
		}
	}
}
