package server

import (
	"bytes"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"strconv"

	"github.com/hashicorp/yamux"
	_ "github.com/mattn/go-sqlite3" // Needed to use sqlite
)

// Start takes a name, seed, and port and starts the universe server
func Start(name string, seed, port int) {
	dbName := name + ".db"

	db, err := sql.Open("sqlite3", dbName)
	checkErr(err)
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS setting (name TEXT PRIMARY KEY, val TEXT)")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS chunk (planet INT, lat INT, lon INT, alt INT, data BLOB, PRIMARY KEY (planet, lat, lon, alt))")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS planet (planet INT PRIMARY KEY, data BLOB)")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS user (name TEXT PRIMARY KEY, data BLOB)")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	rows, err := db.Query("SELECT val FROM setting WHERE name = \"seed\"")
	checkErr(err)
	var val string
	if rows.Next() {
		err = rows.Scan(&val)
		checkErr(err)
		seed, err = strconv.Atoi(val)
	} else {
		stmt, err = db.Prepare("INSERT INTO setting VALUES (\"seed\",?)")
		checkErr(err)
		_, err = stmt.Exec(seed)
		checkErr(err)
	}
	rows.Close()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	u := NewUniverse(0)
	p := NewPlanet(u, 10.0, 1.0, 85.0, 80, 60, 5)

	err = enc.Encode(p.planetState)
	checkErr(err)

	stmt, err = db.Prepare("DROP TABLE IF EXISTS temp")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS temp (data BLOB)")
	checkErr(err)
	_, err = stmt.Exec()
	checkErr(err)

	stmt, err = db.Prepare("INSERT INTO temp VALUES (?)")
	checkErr(err)
	_, err = stmt.Exec(buf.Bytes())
	checkErr(err)

	rows, err = db.Query("SELECT data FROM temp")
	checkErr(err)
	var data []byte
	if rows.Next() {
		err = rows.Scan(&data)
		checkErr(err)
	}
	rows.Close()

	var dbuf bytes.Buffer
	dbuf.Write(data)
	dec := gob.NewDecoder(&dbuf)
	var c planetState
	err = dec.Decode(&c)
	checkErr(err)

	u = NewUniverse(seed)
	p = NewPlanet(u, 10.0, 1.0, 85.0, 80, 64, 16)

	arith := new(Arith)
	listener, e := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	log.Printf("Server listening on port %v...\n", port)
	for {
		conn, e := listener.Accept()
		if e != nil {
			panic(e)
		}

		// Set up server side of yamux
		mux, e := yamux.Server(conn, nil)
		if e != nil {
			panic(e)
		}
		muxConn, e := mux.Accept()
		if e != nil {
			panic(e)
		}
		s := rpc.NewServer()
		s.Register(arith)
		go s.ServeConn(muxConn)

		// Set up stream back to client
		stream, e := mux.Open()
		if e != nil {
			panic(e)
		}
		crpc := rpc.NewClient(stream)

		// Synchronous call
		args := &Args{A: 7, B: 8}
		var reply int
		err = crpc.Call("Arith.Multiply", args, &reply)
		if err != nil {
			log.Fatal("arith error:", err)
		}
		fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)
	}
}
