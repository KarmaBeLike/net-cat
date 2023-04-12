package main

import (
	"fmt"
	"log"
	"net-cat/funcs"
)

func main() {
	s, err := funcs.NewServer("tcp", ":8000", 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on the port: " + "localhost:8000")

	go s.Broadcaster()
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go s.HandleConn(conn)

	}
}
