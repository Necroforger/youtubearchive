package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/jinzhu/gorm"
)

func execSQL(db *gorm.DB, query string) {
	rows, err := db.Raw(query).Rows()
	if err != nil {
		log.Fatal("error executing query: ", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatal(err)
	}

	columnFormatted := strings.Join(columns, "\t")
	fmt.Println(columnFormatted)
	fmt.Println(strings.Repeat("-", len(columnFormatted)+4*len(columns)))

	// Iterate over every row and print with tab separated columns
	for rows.Next() {
		vals := make([]interface{}, len(columns))
		for i := range columns {
			vals[i] = new(sql.RawBytes)
		}

		err := rows.Scan(vals...)
		if err != nil {
			log.Fatal(err)
		}

		for _, v := range vals {
			fmt.Printf("%s\t", []byte(*v.(*sql.RawBytes)))
		}

		fmt.Println()
	}
}
