package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/huanlin/go-htmx-websockets-example/internal/hardware"
)

type server struct {
	subscriberMessageBuffer int
	mux                     http.ServeMux
	subscribersMutex        sync.Mutex
	subscribers             map[*subscriber]struct{}
}

type subscriber struct {
	msgs chan []byte
}

func NewServer() *server {
	s := &server{
		subscriberMessageBuffer: 10,
		subscribers:             make(map[*subscriber]struct{}),
	}

	s.mux.Handle("/", http.FileServer(http.Dir("./htmx")))
	s.mux.HandleFunc("/ws", s.subscribeHandler)
	return s
}

func (s *server) subscribeHandler(writer http.ResponseWriter, req *http.Request) {
	err := s.subscribe(req.Context(), writer, req)
	if err != nil {
		fmt.Println(err)
	}
}

func (s *server) addSubscriber(sub *subscriber) {
	s.subscribersMutex.Lock()
	s.subscribers[sub] = struct{}{}
	s.subscribersMutex.Unlock()
	fmt.Println("Added subscriber", sub)
}

func (s *server) subscribe(ctx context.Context, writer http.ResponseWriter, req *http.Request) error {
	var wsConn *websocket.Conn
	sub := &subscriber{
		msgs: make(chan []byte, s.subscriberMessageBuffer),
	}
	s.addSubscriber(sub)

	wsConn, err := websocket.Accept(writer, req, nil)
	if err != nil {
		return err
	}
	defer wsConn.CloseNow()

	ctx = wsConn.CloseRead(ctx)

	for {
		select {
		case msg := <-sub.msgs:
			ctx, cancel := context.WithTimeout(ctx, time.Second)
			defer cancel()
			err := wsConn.Write(ctx, websocket.MessageText, msg)
			if err != nil {
				return err
			}
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *server) broadcast(msg []byte) {
	s.subscribersMutex.Lock()
	for sub := range s.subscribers {
		sub.msgs <- msg
	}
	s.subscribersMutex.Unlock()
}

func main() {
	fmt.Println("Starting system monitor...")
	srv := NewServer()
	go func(s *server) {
		for {
			systemSection, err := hardware.GetSystemSection()
			if err != nil {
				fmt.Println(err)
			}

			diskSection, err := hardware.GetDiskSection()
			if err != nil {
				fmt.Println(err)
			}
			cpuSection, err := hardware.GetCpuSection()
			if err != nil {
				fmt.Println(err)
			}

			//fmt.Println(systemSection)
			//fmt.Println(diskSection)
			//fmt.Println(cpuSection)

			timeStamp := time.Now().Format("2006-01-02 15:04:05")
			html := `
			<div hx-swap-oob="innerHTML:#update-timestamp"> ` + timeStamp + `</div>
			<div hx-swap-oob="innerHTML:#system-data"> ` + systemSection + `</div>
			<div hx-swap-oob="innerHTML:#disk-data"> ` + diskSection + `</div>
			<div hx-swap-oob="innerHTML:#cpu-data"> ` + cpuSection + `</div>
			`

			s.broadcast([]byte(html))

			time.Sleep(3 * time.Second)
		}
	}(srv)

	err := http.ListenAndServe(":8080", &srv.mux)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
