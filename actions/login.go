package actions

import (
	"fmt"
	"net/url"
	"remote"
	"web/websocket/client"
)

func Login(cl *client.Client, msg *client.Message) {
	qu := url.Values{}
	qu.Set("entity", "User")
	qu.Set("limits", "1")
	qu.Set("fields", "Id")

	user := struct {
		Id       int
		Email    string
		Password string
	}{}
	err := convert(msg.Content, &user)
	if err != nil {
		cl.Io.Errc <- err
	} else {

		qu.Set("Email", user.Email)
		qu.Set("Password", user.Password)
		_, body, err := cl.Rt.GET(remote.RequestConfig{
			URL:   "http://localhost:1234/",
			Query: qu,
		})

		if err != nil {
			cl.Io.Errc <- err
		} else {

			userID := new(struct {
				ID int `json:"Id"`
			})
			err = convertResponseTo(body, userID)
			if err != nil {
				cl.Io.Errc <- err
			} else {
				fmt.Println("login:", userID)

				if err == nil {
					if userID != nil {
						cl.Id = userID.ID
						cl.Store.Clients[cl.Id] = cl
					}
				}
			}
		}
	}
}
