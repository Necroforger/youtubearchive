/*
package main
ytatool performs various common database management functions for youtube-archive

usage: ytatool --database=DATABASE [<flags>] <command> [<args> ...]

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

  exec <sql>
    execute sql and print the results

  get-terminated
    return a list of terminated channels and their channel URLs

  get-active
    return a list of active channels and their channel URLs
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
	database              = kingpin.Flag("database", "the path to the sqlite database to use when performing commands").Short('d').Required().String()
	updateTerminated      = kingpin.Command("update-terminated", "updates the table of terminated channels in the database")
	updateTerminatedProcs = updateTerminated.Flag("procs", "number of http processes to execute concurrently").Short('p').Default("10").Int()

	execCmd    = kingpin.Command("exec", "execute sql and print the results")
	execCmdSQL = execCmd.Arg("sql", "sql string to execute on the database").Required().String()

	terminatedCmd = kingpin.Command("get-terminated", "return a list of terminated channels and their channel URLs")
	activeCmd     = kingpin.Command("get-active", "return a list of active channels and their channel URLs")
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
	case "exec":
		execSQL(db, *execCmdSQL)
	case "get-terminated":
		execSQL(db, "select uploader, uploader_url from terminated_channels where terminated = 1;")
	case "get-active":
		execSQL(db, "select uploader, uploader_url from terminated_channels where terminated = 0;")
	}
	defer db.Close()
}
