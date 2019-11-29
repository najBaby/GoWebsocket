package client

type Kind string

const (
	Run Kind = "run"
)

type Message struct {
	Id      int `json:"-"`
	Kind    Kind
	Content interface{}
}
