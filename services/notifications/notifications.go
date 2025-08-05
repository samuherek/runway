package notifications

import "encoding/json"

type Notifications struct {
	Msgs    []string
	ErrMsgs []string
}

func NewNotifications() *Notifications {
	return &Notifications{
		Msgs:    []string{},
		ErrMsgs: []string{},
	}
}

func (n *Notifications) AddError(value string) {
	n.ErrMsgs = append(n.ErrMsgs, value)
}

func (n *Notifications) AddUnexpectedError() {
	n.ErrMsgs = append(n.ErrMsgs, "Unexpected error occured.")
}

func (n *Notifications) AddMessage(value string) {
	n.Msgs = append(n.Msgs, value)
}

func (n *Notifications) JsonErrors() ([]byte, error) {
	return json.Marshal(n.ErrMsgs)
}

func (n *Notifications) JsonMsgs() ([]byte, error) {
	return json.Marshal(n.Msgs)
}
