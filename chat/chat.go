package chat

type Chat struct {
	started bool
	finished bool
}

func NewChat() *Chat {
	return &Chat{}
}