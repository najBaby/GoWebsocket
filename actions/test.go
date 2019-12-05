package actions

import (
	"remote"
	"web/websocket/client"
)

func Test(cl *client.Client, msg *client.Message) {
	compo := struct {
		Compo      string
		Program    string
		LanguageID int
	}{}
	err := convert(msg.Content, &compo)
	if err != nil {
		cl.Io.Errc <- err
	} else {

		_, body, err := cl.Remote.POST(remote.RequestConfig{
			URL: "",
			Body: struct {
				SourceCode string `json:"source_id"`
				LanguageID int    `json:"language_id"`
			}{
				SourceCode: compo.Program,
				LanguageID: compo.LanguageID,
			},
		})
		if err != nil {
			cl.Io.Errc <- err
		} else {
			token := struct {
				Token string `json:"token"`
			}{}
			err := convert(body, &token)
			if err != nil {
				cl.Io.Errc <- err
			} else {
				_, b, err := cl.Remote.GET(remote.RequestConfig{
					URL: "" + token.Token + "",
				})
				if err != nil {
					cl.Io.Errc <- err
				} else {
					m := new(client.Message)
					m.Kind = "test"
					m.Content = b
					cl.Io.Out <- m
				}
			}
		}
	}
}
