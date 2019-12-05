package actions

import (
	"fmt"
	"remote"
	"time"
	"web/websocket/client"
)

func Composer(cl *client.Client, msg *client.Message) {
	compo := struct {
		Compo string
	}{}
	err := convert(msg.Content, &compo)
	fmt.Println(compo, cl.Store.Rooms)
	if err != nil {
		cl.Io.Errc <- err
	} else {
		if room := cl.Store.GetRoom(compo.Compo); room != nil {
			if room.Admin == cl.ID {
				_, body, err := cl.Remote.POST(remote.RequestConfig{
					URL: "http://localhost:1234/",
					Body: map[string]interface{}{
						"entity": "Challenge",
						"fields": map[string]interface{}{
							"Admin":   room.Admin,
							"Expired": false,
						},
					},
				})
				chalID := struct {
					ID int
				}{}
				convertResponseTo(body, &chalID)
				for _, c := range room.Clients {
					_, _, _ = cl.Remote.POST(remote.RequestConfig{
						URL: "http://localhost:1234/",
						Body: map[string]interface{}{
							"nm2m":     "Challenge",
							"rel":      "Students",
							"nreverse": "Student",
							"m2m": map[string]interface{}{
								"ID": chalID.ID,
							},
							"reverse": map[string]interface{}{
								"ID": c.ID,
							},
						},
					})
				}
				_, _ = body, err
				for i := 10; 0 <= i; i-- {
					<-time.Tick(time.Second)
					m := new(client.Message)
					m.Kind = "ready"
					m.Content = struct {
						T int
					}{
						T: i,
					}
					for _, c := range room.Clients {
						c.Io.Out <- m
					}
				}
				go func() {
					started, _ := time.Parse("15:04:05", "00:00:20")
					for range time.Tick(time.Second) {
						current := started.Format("15:04:05")
						m := new(client.Message)
						m.Kind = "time"
						m.Content = struct {
							Time string
						}{
							Time: started.Format("15:04:05"),
						}
						for _, c := range room.Clients {
							c.Io.Out <- m
						}
						if current == "00:00:00" {
							break
						}
						started = started.Add(-1 * time.Second)
					}
					_, _, _ = cl.Remote.PUT(remote.RequestConfig{
						URL: "http://localhost:1234/",
						Body: map[string]interface{}{
							"entity": "Challenge",
							"filter": map[string]interface{}{
								"ID": chalID.ID,
							},
							"fields": map[string]interface{}{
								"Expired": true,
							},
						},
					})
				}()
			}
		}
	}
}
