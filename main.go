package main

import (
	"fmt"
	"log"
	"net-cat/funcs"
	"os"
)

func main() {
	args := os.Args[1:]
	var port string
	if len(args) != 1 && len(args) != 0 {
		log.Fatal("try again")
	} else if len(args) == 1 {
		port = args[0]
	} else {
		port = "8989"
	}
	s, err := funcs.NewServer("tcp", ":"+port, 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Listening on the port: " + "localhost:" + port + "\n")

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
