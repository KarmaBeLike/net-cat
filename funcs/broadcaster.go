package funcs

import (
	"fmt"
	"net"
	"strings"
	"time"
)

type client chan<- string // Канал исходящих сообщений
var (
	clients  = make(map[string]net.Conn) // Все подключенные клиенты
	entering = make(chan message)
	history  = []string{}
	leaving  = make(chan message)
	messages = make(chan message) // Все входящие сообщения клиента
)

const timeFormat = "2006-01-02 15:04:05"

func (s *server) Broadcaster() {
	for {
		select {

		case msg := <-messages: // Широковещательное входящее сообщение во все каналы исходящих сообщений для клиентов
			s.mutex.Lock()
			for cli, conn := range s.clients {
				if cli != msg.name {
					fmt.Fprintf(conn, "\r%v\r", strings.Repeat(" ", len(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), cli))))
					conn.Write([]byte(fmt.Sprintf("[%s][%s]:%s", msg.time, msg.name, msg.text)))
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
				}
			}
			s.mutex.Unlock()
		case msg := <-entering:
			s.mutex.Lock()
			for cli, conn := range s.clients {
				if cli == msg.name {
					for _, v := range history {
						conn.Write([]byte(fmt.Sprintf("%s", v)))
						// conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
					}
					conn.Write([]byte(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), msg.name)))
				}
				if cli != msg.name {
					fmt.Fprintf(conn, "\r%v\r", strings.Repeat(" ", len(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), cli))))
					conn.Write([]byte(fmt.Sprintf("%s", msg.text)))
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
				}
			}
			s.mutex.Unlock()
		case msg := <-leaving:
			s.mutex.Lock()
			for cli, conn := range s.clients {
				if cli != msg.name {
					fmt.Fprintf(conn, "\r%v\r", strings.Repeat(" ", len(fmt.Sprintf("[%s][%s]:", time.Now().Format(timeFormat), cli))))
					conn.Write([]byte(fmt.Sprintf("%s", msg.text)))
					conn.Write([]byte(fmt.Sprintf("\n[%s][%s]:", time.Now().Format(timeFormat), cli)))
				}
			}
			s.mutex.Unlock()
		}
	}
}
