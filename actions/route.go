package actions

import "web/websocket/client"

type route struct {
	path   client.Kind
	action func(*client.Client, *client.Message)
}

func NewHandlers(routes []route) *map[client.Kind]func(*client.Client, *client.Message) {
	handlers := make(map[client.Kind]func(*client.Client, *client.Message))
	for _, route := range routes {
		handlers[route.path] = route.action
	}
	return &handlers
}
