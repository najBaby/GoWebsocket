package actions

import (
	"fmt"
	"net/url"
	"remote"
	"web/websocket/client"
)

func Alllanguages(cl *client.Client, msg *client.Message) {
	qu := url.Values{}
	qu.Set("entity", "Language")
	qu.Set("limit", "0")
	qu.Set("fields", "Name,ID")

	_, body, err := cl.Remote.GET(remote.RequestConfig{
		URL:   "http://localhost:1234/",
		Query: qu,
	})

	if err != nil {
		cl.Io.Errc <- err
	} else {
		languages := []*struct {
			Name  string
			ID int
		}{}
		err := convertResponseTo(body, &languages)
		fmt.Println(body, languages)
		if err != nil {
			cl.Io.Errc <- err
		} else {
			m := new(client.Message)
			m.ID = cl.ID
			m.Kind = "alllanguages"
			m.Content = languages
			cl.Io.Out <- m
		}
	}
}
