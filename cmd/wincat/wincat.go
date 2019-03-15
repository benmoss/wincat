package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	port := os.Args[1]
	if port == "" {
		log.Fatal("usage: wincat [port]")
	}
	conn, err := net.Dial("tcp", "localhost:"+port)
	if err != nil {
		log.Fatal(fmt.Sprintf("Failed to connect to localhost:%s because %s", port, err))
	}
	go func() {
		_, e := io.Copy(conn, os.Stdin)
		if e != nil {
			log.Fatal("Attempted to copy from stdin, but got error:", e)
		}
	}()
	_, e := io.Copy(os.Stdout, conn)
	if e != nil {
		log.Fatal("Attempted to copy to stdout, but got error:", e)
	}
}
