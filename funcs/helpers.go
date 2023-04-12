package funcs

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type server struct {
	Listener       net.Listener
	clients        map[string]net.Conn
	maxConnections int
	connections    int
	mutex          sync.Mutex
}

type message struct {
	time string
	name string
	text string
}

func NewServer(protocol, addres string, maxConnections int) (*server, error) {
	listener, err := net.Listen(protocol, addres)
	if err != nil {
		return nil, err
	}
	return &server{
		Listener:       listener,
		clients:        make(map[string]net.Conn),
		maxConnections: maxConnections,
		mutex:          sync.Mutex{},
	}, nil
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}

func Welcome(conn net.Conn) {
	file, err := os.ReadFile("logo.txt")
	if err != nil {
		log.Print("File is not correct\n")
		return
	}
	welcometext := string(file)
	fmt.Fprintf(conn, welcometext)
}
