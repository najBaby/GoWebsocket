package actions

import (
	"net/url"
	"remote"
	"web/websocket/client"
)

func Allusers(cl *client.Client, msg *client.Message) {
	qu := url.Values{}
	qu.Set("entity", "User")
	qu.Set("limits", "0")
	qu.Set("fields", "Id,Name,Language")

	_, body, err := cl.Rt.GET(remote.RequestConfig{
		URL:   "http://localhost:1234/",
		Query: qu,
	})

	if err != nil {
		cl.Io.Errc <- err
	} else {

		users := []*struct {
			Id       int
			Name     string
			Language string
		}{}
		err := convertResponseTo(body, &users)
		if err != nil {
			cl.Io.Errc <- err
		} else {

			for _, user := range users {
				if user.Id == cl.Id {

					break
				}
			}
			m := new(client.Message)
			m.Id = cl.Id
			m.Kind = "allusers"
			m.Content = users
			cl.Io.Out <- m
		}
	}
}
