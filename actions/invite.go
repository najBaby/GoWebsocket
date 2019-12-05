package actions

import (
	"crypto/md5"
	"encoding/hex"
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
				Name     string
				Language struct {
					Name string
				}
			}{}
			err := convertResponseTo(body, &user)
			fmt.Println(user)
			if err != nil {
				cl.Io.Errc <- err
			} else {
				put := map[string]interface{}{
					"entity": "Student",
					"filter": map[string]interface{}{
						"ID": cl.ID,
					},
					"fields": map[string]interface{}{
						"ID":    cl.ID,
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
					room := client.NewRoom()
					room.Name = hex.EncodeToString(md5.New().Sum([]byte(strconv.Itoa(cl.ID))))
					room.Admin = cl.ID
					room.AddClient(cl)
					cl.Store.AddRoom(room)

					m := new(client.Message)
					m.Kind = "invite"
					m.Content = struct {
						User struct {
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

					fmt.Println(room.Name)
					for _, id := range receivers.Receivers {
						if c, ok := cl.Store.Clients[id]; ok {
							c.Io.Out <- m
						}
					}
					m = new(client.Message)
					m.Kind = "busy"
					m.Content = map[string]interface{}{
						"Busy": room.Name,
					}
					cl.Io.Out <- m
				}
			}
		}
	}
}
