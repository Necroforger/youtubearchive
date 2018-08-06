/*
youtube-archive-server serves a webserver
that reads from an archived database of youtube videos
*/
package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/chi"

	"github.com/jinzhu/gorm"

	"github.com/Necroforger/youtubearchive/server"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	dbpath = flag.String("db", "", "path to the database sqlite3 file")
	addr   = flag.String("addr", ":80", "address to bind to")
	tpldir = flag.String("templates", "./templates/*", "directory containing site templates")
	static = flag.String("static", "./static", "directory containing static files")
)

func main() {
	flag.Parse()

	if *dbpath == "" {
		log.Fatal("please supply a database path with -db")
	}

	db, err := gorm.Open("sqlite3", *dbpath)
	if err != nil {
		log.Fatal("Err opening database: ", err)
	}
	defer db.Close()

	s := server.NewServer(db)
	err = s.LoadTemplatesGlob(*tpldir)
	if err != nil {
		log.Fatal(err)
	}
	s.Logger = os.Stderr

	r := chi.NewMux()
	r.Use(
		middleware.Logger,
		middleware.RealIP,
	)
	r.Mount("/", s.GetRoutes())

	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir(*static))))
	http.Handle("/", r)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
