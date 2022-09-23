package chat

import "fmt"

type Chat struct {
	started  bool
	finished bool
	Me       *User
	Message  []*Message
}

type Message struct {
	ID int32
	NickName string
	Message string
}

func NewChat() *Chat {
	return &Chat{}
}

func (c *Chat) Talk(m string, u *User) (bool, error) {
	if c.finished {
		return true, nil
	}

	c.Message = append(c.Message, &Message{
		ID:       u.ID,
		NickName: u.NickName,
		Message:  m,
	})
	c.Display(c.Me)

	return false, nil
}

func (c *Chat) Display(me *User) {
	fmt.Println("======================")
	for _, v := range c.Message {
		fmt.Println(v)
		fmt.Println("- - - - - - - - - - - -")
	}
	fmt.Println("======================")
}
