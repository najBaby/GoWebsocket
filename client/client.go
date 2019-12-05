package client

import (
	"encoding/json"
	"fmt"
	"remote"

	"golang.org/x/net/websocket"
)

type io struct {
	In   chan *Message
	Out  chan *Message
	Errc chan error
}

type Client struct {
	Io     io
	ID     int
	Msg    Message
	Store  *store
	Remote *remote.Remote
}

func newClient(store *store) *Client {
	return &Client{
		Io: io{
			In:   make(chan *Message),
			Out:  make(chan *Message),
			Errc: make(chan error, 1),
		},
		Store:  store,
		Remote: remote.NewRemote(remote.RemoteConfig{}),
	}
}

func (cl *Client) handle(conn *websocket.Conn, handlers *map[Kind]func(*Client, *Message)) {
	go func() {
		dec := json.NewDecoder(conn)
		for {
			var m Message
			if err := dec.Decode(&m); err != nil {
				cl.Io.Errc <- err
				return
			}
			cl.Io.In <- &m

		}
	}()

	go func() {
		enc := json.NewEncoder(conn)
		for m := range cl.Io.Out {
			if err := enc.Encode(m); err != nil {
				cl.Io.Errc <- err
				return
			}
		}
	}()

	for {
		select {
		case msg := <-cl.Io.In:
			h := *handlers
			if f, ok := h[msg.Kind]; ok {
				f(cl, msg)
			}
		case errc := <-cl.Io.Errc:
			fmt.Println("errc:", errc)
		}
	}
}
