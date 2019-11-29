package actions

import (
	"fmt"
	"net/url"
	"remote"
	"strconv"
	"web/websocket/client"
)

func Invite(cl *client.Client, msg *client.Message) {
	receivers := struct {
		Receivers []int
	}{}
	err := convert(msg.Content, &receivers)
	fmt.Println(receivers)

	if err != nil {
		cl.Io.Errc <- err
	} else {

		qu := url.Values{}
		qu.Set("entity", "User")
		qu.Set("Id", strconv.Itoa(cl.Id))
		qu.Set("limits", "1")
		qu.Set("fields", "Name,Language")

		_, body, err := cl.Rt.GET(remote.RequestConfig{
			URL:   "http://localhost:1234/",
			Query: qu,
		})
		if err != nil {
			cl.Io.Errc <- err
		} else {

			user := struct {
				Name     string
				Language string
			}{}
			err := convertResponseTo(body, &user)

			if err != nil {
				cl.Io.Errc <- err
			} else {
				room := client.NewRoom()
				room.Name = "compo1"
				room.Admin = cl.Id
				room.AddClient(cl)
				cl.Store.AddRoom(room)

				m := new(client.Message)
				m.Id = cl.Id
				m.Kind = "invite"
				m.Content = struct {
					User struct {
						Name     string
						Language string
					}
					Compo string
				}{
					User:  user,
					Compo: room.Name,
				}

				fmt.Println(cl.Store.Clients)
				for _, id := range receivers.Receivers {
					if c, ok := cl.Store.Clients[id]; ok {
						c.Io.Out <- m
					}
				}
			}
		}
	}

}
