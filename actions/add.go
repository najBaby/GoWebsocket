package actions

import (
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
					Id       int
					Name     string
					Language string
				}{}
				err := convertResponseTo(body, &user)
				if err != nil {
					cl.Io.Errc <- err
				} else {
					m := new(client.Message)
					user.Id = cl.Id
					m.Kind = "add"
					m.Content = struct {
						User struct {
							Id       int
							Name     string
							Language string
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
