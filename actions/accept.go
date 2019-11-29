package actions

import (
	"fmt"
	"web/websocket/client"
)

func Accept(cl *client.Client, msg *client.Message) {
	compo := struct {
		Id    int
		Compo string
	}{}
	err := convert(msg.Content, &compo)
	fmt.Println(compo)
	if err != nil {
		cl.Io.Errc <- err
	} else {
		if room := cl.Store.GetRoom(compo.Compo); room != nil {
			if c := cl.Store.Clients[compo.Id]; c != nil && room.Admin == cl.Id {
				room.AddClient(c)
			}
		}
	}
}
