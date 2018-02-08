package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/lib/pq"
)

// errorf prints to stderr and exits.
func errorf(format string, args ...interface{}) {
	m := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(m, "\n") {
		m += "\n"
	}
	fmt.Fprintf(os.Stderr, m)
	os.Exit(1)
}

// main is the entry point of the program.
func main() {
	// Command line flags.
	delimiter := flag.String("delimiter", ",", "delimeter")
	driver := flag.String("driver", "postgres", "sql driver, defaults to postgres (mssql|postgres)")
	source := flag.String("source", "", "data source name")
	query := flag.String("query", "-", "sql query")
	flag.Parse()

	// Make sure we got a valid parameter
	if *source == "" {
		errorf("missing required -source flag")
	}

	// Make sure that we got a valid driver.
	if *driver != "mssql" && *driver != "postgres" {
		errorf("%v is not a valid driver, postgres is the only supported driver for now", *driver)
	}

	// Connect to the server.
	db, err := sql.Open(*driver, *source)
	if err != nil {
		errorf("unable to open database connection using driver %s and source %s, errror: %v", *driver, *source, err)
	}
	defer db.Close()

	// Ping the database.
	if err = db.Ping(); err != nil {
		errorf("can't ping database: %v", err)
	}

	// Accept query to be sent through stdin as well.
	if *query == "-" {
		q, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			errorf("%v", err)
		}
		*query = string(q)
	}

	// Execute the query.
	rows, err := db.Query(*query)
	if err != nil {
		errorf("coudln't execute query %s, error: %v\n", *query, err)
	}

	// Calculate number of columns for the returned rows.
	cols, err := rows.Columns()
	if err != nil {
		errorf("%v\n", err)
	}

	// Prepare data structures for the row data.
	res := make([][]byte, len(cols))
	dest := make([]interface{}, len(cols))
	for i := range res {
		dest[i] = &res[i]
	}

	// Iterate over each row and print the results.
	for rows.Next() {
		err = rows.Scan(dest...)
		if err != nil {
			errorf("%v", err)
		}

		// Create a temporary string where we store each column data.
		// We append the delimiter after each column.
		var tmp string
		for _, r := range res {
			tmp += string(r) + *delimiter
		}

		// Print the temp string to stdout, we also remove the last delimiter before we print anything.
		fmt.Fprintf(os.Stdout, "%s\n", tmp[0:len(tmp)-1])
	}
}
