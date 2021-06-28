package notification

type Notification struct {
	ClientID  string
	Message   string
	messageID uint64
}

func (n *Notification) SetMessageID(id uint64) {
	n.messageID = id
}
