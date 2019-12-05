package actions

import (
	"fmt"
	"remote"
	"web/websocket/client"
)

func Accept(cl *client.Client, msg *client.Message) {
	compo := struct {
		ID    int
		Compo string
	}{}
	err := convert(msg.Content, &compo)
	fmt.Println(compo)
	if err != nil {
		cl.Io.Errc <- err
	} else {
		if room := cl.Store.GetRoom(compo.Compo); room != nil {
			if c := cl.Store.Clients[compo.ID]; c != nil && room.Admin == cl.ID {
				put := map[string]interface{}{
					"entity": "Student",
					"filter": map[string]interface{}{
						"ID": c.ID,
					},
					"fields": map[string]interface{}{
						"ID":    c.ID,
						"State": 0,
					},
				}
				_, _, err = cl.Remote.PUT(remote.RequestConfig{
					URL:  "http://localhost:1234/",
					Body: put,
				})
				if err != nil {
					cl.Io.Errc <- err
				} else {
					room.AddClient(c)
					m := new(client.Message)
					m.Kind = "busy"
					m.Content = map[string]interface{}{
						"Busy": room.Name,
					}
					c.Io.Out <- m
				}
			}
		}
	}
}
