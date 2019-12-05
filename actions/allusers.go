package actions

import (
	"fmt"
	"net/url"
	"remote"
	"web/websocket/client"
)

func Allusers(cl *client.Client, msg *client.Message) {
	qu := url.Values{}
	qu.Set("entity", "Student")
	qu.Set("limit", "0")
	qu.Set("fields", "ID,Name,Language__Name")

	_, body, err := cl.Remote.GET(remote.RequestConfig{
		URL:   "http://localhost:1234/",
		Query: qu,
	})

	if err != nil {
		cl.Io.Errc <- err
	} else {

		users := []struct {
			ID       int
			Name     string
			Language struct {
				Name string
			}
		}{}
		err := convertResponseTo(body, &users)
		fmt.Println(body, users)
		if err != nil {
			cl.Io.Errc <- err
		} else {

			for index, user := range users {
				if user.ID == cl.ID {
					users = append(users[:index], users[index+1:]...)
					break
				}
			}

			m := new(client.Message)
			m.ID = cl.ID
			m.Kind = "allusers"
			m.Content = users
			cl.Io.Out <- m
		}
	}
}
