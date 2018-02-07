package main

import (
	"database/sql"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

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
	driver := flag.String("driver", "postgres", "sql server driver (can only be postgres currently)")
	name := flag.String("name", "", "name of the database")
	password := flag.String("password", "", "password")
	port := flag.Int("port", 0, "port")
	server := flag.String("server", "", "address")
	user := flag.String("user", "", "user")
	query := flag.String("query", "", "sql query")

	// Verify that the required flags are passed.
	seen := make(map[string]bool)
	flag.Parse()
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, r := range []string{"name", "password", "port", "server", "user", "query"} {
		if !seen[r] {
			errorf("missing required -%s flag", r)
		}
	}

	// Make sure that we got a valid driver.
	if *driver != "postgres" {
		errorf("%v is not a valid driver, postgres is the only supported driver for now", *driver)
	}

	// Connect to the server.
	cstr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", *server, *port, *user, *password, *name)
	db, err := sql.Open(*driver, cstr)
	if err != nil {
		errorf("%v", err)
	}
	defer db.Close()

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
		errorf("%v\n", err)
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
