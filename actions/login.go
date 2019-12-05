package actions

import (
	"fmt"
	"net/url"
	"remote"
	"web/websocket/client"
)

func Login(cl *client.Client, msg *client.Message) {
	qu := url.Values{}
	qu.Set("entity", "Student")
	qu.Set("limit", "1")
	qu.Set("fields", "ID,Name,Language__Name,Language__Index")

	user := struct {
		Email    string
		Password string
	}{}
	err := convert(msg.Content, &user)
	if err != nil {
		cl.Io.Errc <- err
	} else {

		qu.Set("Email", user.Email)
		qu.Set("Password", user.Password)
		_, body, err := cl.Remote.GET(remote.RequestConfig{
			URL:   "http://localhost:1234/",
			Query: qu,
		})

		if err != nil {
			cl.Io.Errc <- err
		} else {

			userID := new(struct {
				ID       int
				Name     string
				Language struct {
					Name  string
					Index int
				}
			})
			err = convertResponseTo(body, userID)
			if err != nil {
				cl.Io.Errc <- err
			} else {
				fmt.Println("login:", userID)

				if err == nil {
					if userID != nil {
						cl.ID = userID.ID
						cl.Store.Clients[cl.ID] = cl
						m := new(client.Message)
						m.Kind = "connected"
						m.Content = map[string]interface{}{
							"Name": userID.Name,
							"Language": map[string]interface{}{
								"Name":  userID.Language.Name,
								"Index": userID.Language.Index,
							},
						}
						cl.Io.Out <- m
						fmt.Println(m)
					}
				}
			}
		}
	}
}
