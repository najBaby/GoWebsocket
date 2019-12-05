package client

type Kind string

type Message struct {
	ID      int `json:"-"`
	Kind    Kind
	Content interface{}
}
