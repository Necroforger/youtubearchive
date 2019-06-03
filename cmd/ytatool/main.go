/*
package main
ytatool performs various common database management functions for youtube-archive

usage: ytatool [<flags>] <command> [<args> ...]

Flags:
      --help               Show context-sensitive help (also try --help-long and
                           --help-man).
  -d, --database=DATABASE  the path to the sqlite database to use when
                           performing commands

Commands:
  help [<command>...]
    Show help.


  update-terminated [<flags>]
    updates the table of terminated channels in the database

    -p, --procs=10  number of http processes to execute concurrently
*/

package main

import (
	"log"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/alecthomas/kingpin"
	"github.com/jinzhu/gorm"
)

var (
	database              = kingpin.Flag("database", "the path to the sqlite database to use when performing commands").Short('d').String()
	updateTerminated      = kingpin.Command("update-terminated", "updates the table of terminated channels in the database")
	updateTerminatedProcs = updateTerminated.Flag("procs", "number of http processes to execute concurrently").Short('p').Default("10").Int()
)

func openDatabase() *gorm.DB {
	if _, err := os.Stat(*database); err == os.ErrNotExist {
		log.Fatal("database file not found")
	}

	db, err := gorm.Open("sqlite3", *database)
	if err != nil {
		log.Fatal("error opening database: ", err)
	}

	return db
}

func main() {
	cmd := kingpin.Parse()

	db := openDatabase()
	switch cmd {
	case "update-terminated":
		updateTerminatedCmd(db)
	}
	defer db.Close()
}
