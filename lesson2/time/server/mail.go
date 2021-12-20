package main

import (
	"bufio"
	"context"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Server struct {
	listener    net.Listener
	Connections chan net.Conn
}

func NewServer(address string) Server {
	lister, err := net.Listen("tcp", ":8001")
	if err != nil {
		panic(err)
	}

	connChan := make(chan net.Conn)

	return Server{
		listener:    lister,
		Connections: connChan,
	}
}

func (s Server) Start() {
	log.Printf("server started on %s\n", s.listener.Addr())

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		s.Connections <- conn
	}
}

func pipeMessage(ctx context.Context, messageChan chan string) {
	select {
	case <-ctx.Done():
		log.Println("Stop messenger")
		return
	default:
		for {
			sc := bufio.NewScanner(os.Stdin)
			sc.Scan()
			messageChan <- sc.Text() + "\n"
		}
	}
}

func sendToChan(chans []chan string, message string) {
	for i := range chans {
		chans[i] <- message
	}
}

func broker(ctx context.Context, chanChan chan chan string, messageMain chan string) {
	chans := make([]chan string, 0)
	for {
		select {
		case <-ctx.Done():
			return
		case message := <-messageMain:
			sendToChan(chans, message)
		case newChan := <-chanChan:
			chans = append(chans, newChan)
		}
	}
}

func main() {
	srv := NewServer(":8001")
	go srv.Start()

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	wg := sync.WaitGroup{}
	messageMain := make(chan string)
	channels := make(chan chan string)
	go pipeMessage(ctx, messageMain)
	go broker(ctx, channels, messageMain)
	for {
		select {
		case <-ctx.Done():
			log.Println("start graceful")
			wg.Wait()
			log.Println("stop graceful")
			return
		case conn := <-srv.Connections:
			wg.Add(1)
			messageChan := make(chan string)
			channels <- messageChan
			go handleConn(ctx, conn, &wg, messageChan)
		}
	}
}

func handleConn(ctx context.Context, c net.Conn, wg *sync.WaitGroup, messageChan chan string) {
	defer func() {
		wg.Done()
		c.Close()
	}()

	for {
		t := time.NewTicker(time.Second)

		select {
		case <-ctx.Done():
			io.WriteString(c, "Bye!")
			return
		case now := <-t.C:
			io.WriteString(c, now.Format("15:04:05\n\r"))
		case message := <-messageChan:
			io.WriteString(c, message)
		}
	}
}
