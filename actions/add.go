package actions

import (
	"fmt"
	"net/url"
	"remote"
	"strconv"
	"web/websocket/client"
)

func Add(cl *client.Client, msg *client.Message) {
	compo := struct {
		Compo string
	}{}
	err := convert(msg.Content, &compo)
	if err != nil {
		cl.Io.Errc <- err
	} else {
		if room := cl.Store.GetRoom(compo.Compo); room != nil {
			fmt.Println("ADD", cl.ID)
			qu := url.Values{}
			qu.Set("entity", "Student")
			qu.Set("ID", strconv.Itoa(cl.ID))
			qu.Set("limit", "1")
			qu.Set("fields", "Name,Language__Name")

			_, body, err := cl.Remote.GET(remote.RequestConfig{
				URL:   "http://localhost:1234/",
				Query: qu,
			})
			if err != nil {
				cl.Io.Errc <- err
			} else {

				user := struct {
					ID       int
					Name     string
					Language struct {
						Name string
					}
				}{}
				err := convertResponseTo(body, &user)
				if err != nil {
					cl.Io.Errc <- err
				} else {
					m := new(client.Message)
					user.ID = cl.ID
					m.Kind = "add"
					m.Content = struct {
						User struct {
							ID       int
							Name     string
							Language struct {
								Name string
							}
						}
						Compo string
					}{
						User:  user,
						Compo: room.Name,
					}
					if admin := room.GetClient(room.Admin); admin != nil {
						admin.Io.Out <- m
					}
				}
			}
		}
	}
}
