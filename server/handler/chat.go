package handler

import (
	"chat/build"
	"chat/chat"
	"chat/gen/pb"
	"fmt"
	"sync"
)

type Chat struct {
	sync.RWMutex
	chats  map[int32]*chat.Chat
	client map[int32][]pb.ChatService_ChatServer
}

func NewChat() *Chat {
	return &Chat{
		chats:  make(map[int32]*chat.Chat),
		client: make(map[int32][]pb.ChatService_ChatServer),
	}
}

func (c *Chat) Chat(stream pb.ChatService_ChatServer) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}

		roomID := req.GetRoomId()
		user := build.User(req.GetUser())

		switch req.GetAction().(type) {
		// チャット開始リクエスト
		case *pb.ChatRequest_Start:
			err := c.start(stream, roomID, user)
			if err != nil {
				return err
			}
		// チャット
		case *pb.ChatRequest_Talk:
			msg := req.GetTalk().GetMessage()
			err := c.chat(roomID, msg, user)
			if err != nil {
				return err
			}
		}
	}
}

func (c *Chat) start(stream pb.ChatService_ChatServer, roomID int32, me *chat.User) error {
	c.Lock()
	defer c.Unlock()

	ch := c.chats[roomID]
	if ch == nil {
		c.chats[roomID] = chat.NewChat()
		c.client[roomID] = make([]pb.ChatService_ChatServer, 0, 2)
	}
	c.client[roomID] = append(c.client[roomID], stream)

	if len(c.client[roomID]) == 2 {
		for _, s := range c.client[roomID] {
			err := s.Send(&pb.ChatResponse{
				Event: &pb.ChatResponse_Ready{
					Ready: &pb.ChatResponse_ReadyEvent{},
				},
			})
			if err != nil {
				return err
			}
		}
		fmt.Printf("chat has started room_id=%v\n", roomID)
	} else {
		err := stream.Send(&pb.ChatResponse{
			Event: &pb.ChatResponse_Waiting{
				Waiting: &pb.ChatResponse_WaitingEvent{},
			},
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Chat) chat(roomID int32, msg string, u *chat.User) error {
	c.Lock()
	defer c.Unlock()

	ch := c.chats[roomID]
	finished, err := ch.Talk(msg, u)
	if err != nil {
		return err
	}

	for _, v := range c.client[roomID] {
		err := v.Send(&pb.ChatResponse{
			Event: &pb.ChatResponse_Chated{
				Chated: &pb.ChatResponse_ChatedEvent{
					// 一旦仮のチャットログ
					ChatLogs: []string{"hoge, hoge"},
				},
			},
		})
		if err != nil {
			return err
		}

		if finished {
			err := v.Send(&pb.ChatResponse{
				Event: &pb.ChatResponse_Finished{
					Finished: &pb.ChatResponse_FinishedEvent{},
				},
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
