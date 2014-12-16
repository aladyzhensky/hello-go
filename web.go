package main

import (
	"code.google.com/p/log4go"
	"fmt"
	"launchpad.net/goyaml"
	"net/http"
	"os"
	"runtime"
	"database/sql"
      _ "github.com/go-sql-driver/mysql"
        "encoding/json"
)

const (
	HostVar = "VCAP_APP_HOST"
	PortVar = "VCAP_APP_PORT"
)

type T struct {
	A string
	B []int
}

type ClearDBInfo struct {
Credentials ClearDBCredentials `json:"credentials"`
}
 
type ClearDBCredentials struct {
hostname string `json:"host"`
port string `json:"port"`
 string `json:"user"`
Password string `json:"password"`
}
  

func main() {
	log := make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

	http.HandleFunc("/", hello)
	var port string
	if port = os.Getenv(PortVar); port == "" {
		port = "8080"
	}
	log.Debug("Listening at port %v\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
	
	 s := os.Getenv("VCAP_SERVICES")
	services := make(map[string][]ClearDBInfo)
	err := json.Unmarshal([]byte(s), &services)
	if err != nil {
	log.Debug("Error parsing MySQL connection information: %v\n", err.Error())
	return
	}
 
	info := services["cleardb"]
	if len(info) == 0 {
	log.Debug("No ClearDB databases are bound to this application.\n")
	return
	}
 
	// Assumes only a single ClearDB is bound to this application
	creds := info[0].Credentials
 
	host := creds.hostname
	port := creds.port
	name := creds.user
	password := creds.Password
 
// Use host, port, user and password to connect to MySQL using the chosen driver
}

func hello(res http.ResponseWriter, req *http.Request) {
		// Dump Go version
	fmt.Fprintf(res, "%v\n\n", runtime.Version())

	// Dump ENV
	env := os.Environ()
	for _, e := range env {
		fmt.Fprintln(res, e)
	}

	//Dump some YAML
	t := T{A: "Foo", B: []int{1, 2, 3}}
	if d, err := goyaml.Marshal(&t); err != nil {
		fmt.Fprintf(res, "Unable to dump YAML")
	} else {
		fmt.Fprintf(res, "\n\n--- \n%s", d)
	}
	
	db, err := sql.Open("mysql", "user:password@/dbname")
    if err != nil {
        panic(err.Error())  // Just for example purpose. You should use proper error handling instead of panic
    }
    defer db.Close()

    // Execute the query
    rows, err := db.Query("SELECT * FROM table")
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Get column names
    columns, err := rows.Columns()
    if err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }

    // Make a slice for the values
    values := make([]sql.RawBytes, len(columns))

    // rows.Scan wants '[]interface{}' as an argument, so we must copy the
    // references into such a slice
    // See http://code.google.com/p/go-wiki/wiki/InterfaceSlice for details
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }

    // Fetch rows
    for rows.Next() {
        // get RawBytes from data
        err = rows.Scan(scanArgs...)
        if err != nil {
            panic(err.Error()) // proper error handling instead of panic in your app
        }

        // Now do something with the data.
        // Here we just print each column as a string.
        var value string
        for i, col := range values {
            // Here we can check if the value is nil (NULL value)
            if col == nil {
                value = "NULL"
            } else {
                value = string(col)
            }
            fmt.Println(columns[i], ": ", value)
        }
        fmt.Println("-----------------------------------")
    }
    if err = rows.Err(); err != nil {
        panic(err.Error()) // proper error handling instead of panic in your app
    }
}
