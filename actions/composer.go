package actions

import (
	"fmt"
	"strconv"
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
			if room.Admin == cl.Id {
				go func() {
					ticker := time.NewTicker(time.Second)
					defer ticker.Stop()
					done := make(chan bool)
					go func() {
						time.Sleep(10 * time.Second)
						done <- true
					}()
					i := 11
					for {
						select {
						case <-done:
							fmt.Println("Done!")
							return
						case <-ticker.C:
							m := new(client.Message)
							m.Kind = "ready"
							m.Content = struct {
								T string
							}{
								T: strconv.Itoa(i),
							}
							for _, c := range room.Clients {
								c.Io.Out <- m
							}
						}
						i--
					}
				}()
				time.AfterFunc(10*time.Second, func() {
					ticker := time.NewTicker(time.Second)
					defer ticker.Stop()
					done := make(chan bool)
					go func() {
						time.Sleep(20 * time.Second)
						done <- true
					}()
					for {
						select {
						case <-done:
							fmt.Println("Done!")
							return
						case t := <-ticker.C:
							m := new(client.Message)
							m.Kind = "time"
							m.Content = struct {
								T string
							}{
								T: t.Format("15:04:05"),
							}
							for _, c := range room.Clients {
								c.Io.Out <- m
							}
						}
					}
				})
			}
		}
	}
}
