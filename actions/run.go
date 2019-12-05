package actions

import (
	"fmt"
	"net/url"
	"remote"
	"strconv"
	"web/websocket/client"
)

func Run(cl *client.Client, msg *client.Message) {
	compo := struct {
		Compo    string
		Exercise int
		Program  string
		Index    int
	}{}

	err := convert(msg.Content, &compo)
	if err != nil {
		cl.Io.Errc <- err
	} else {
		req := struct {
			Code           string `json:"code"`
			ExecptedResult string `json:"excepted_result"`
			Language       int    `json:"language"`
			Result         string `json:"result"`
			APIKey         string `json:"apikey"`
		}{}
		req.ExecptedResult = "Hello World"
		req.Language = compo.Index
		req.Code = compo.Program
		req.APIKey = "99e86a266c387f47e799ec662db03c3eeed84cb6d8a83433fce60a7013b01924dd688fc328b2459cc793cde7850957819702bbbb58fdae5eb809c91e0c88248f"
		_, body, err := cl.Remote.POST(remote.RequestConfig{
			URL:  "http://192.168.50.67:8000/compile",
			Body: req,
		})

		if err != nil {
			cl.Io.Errc <- err
		} else {

			res := struct {
				Time     float64 `json:"time"`
				Error    bool    `json:"error"`
				Valide   bool    `json:"valide"`
				Status   int     `json:"status"`
				Resultat string  `json:"resultat"`
				Message  string  `json:"message"`
			}{}
			fmt.Println(body)
			err := convert(body, &res)
			if err != nil {
				cl.Io.Errc <- err
			} else {

				if room := cl.Store.GetRoom(compo.Compo); room != nil {
					qu := url.Values{}
					qu.Set("entity", "Student")
					qu.Set("ID", strconv.Itoa(cl.ID))
					qu.Set("limit", "1")
					qu.Set("fields", "Name,Language")

					_, body, err := cl.Remote.GET(remote.RequestConfig{
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
							m := new(client.Message)
							m.ID = cl.ID
							m.Kind = "result"
							m.Content = struct {
								Error    bool
								Resultat string
								Valide   bool
								Message  string
								Time     float64
							}{
								Error:    res.Error,
								Resultat: res.Resultat,
								Valide:   res.Valide,
								Message:  res.Message,
								Time:     res.Time,
							}
							for _, c := range room.Clients {
								c.Io.Out <- m
							}
						}
					}
				}
			}
		}
	}
}
