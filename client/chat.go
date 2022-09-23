package client

import (
	"bufio"
	"chat/build"
	"chat/chat"
	"chat/gen/pb"
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Chat struct {
	sync.RWMutex
	started  bool
	finished bool
	room     *chat.Room
	me       *chat.User
}

func NewChat() *Chat {
	return &Chat{}
}

func (c *Chat) Run() int {
	if err := c.run(); err != nil {
		fmt.Println(err)
		return 1
	}
	return 0
}

func (c *Chat) run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "Failed to connect to grpc server")
	}
	defer conn.Close()

	err = c.matching(ctx, pb.NewMatchServiceClient(conn))
	if err != nil {
		return err
	}

	return c.chat(ctx, pb.NewChatServiceClient(conn))
}

func (c *Chat) matching(ctx context.Context, cli pb.MatchServiceClient) error {
	stream, err := cli.JoinRoom(ctx, &pb.JoinRoomRequest{})
	if err != nil {
		return err
	}
	defer stream.CloseSend()
	fmt.Println("Requested matching...")

	for {
		resp, err := stream.Recv()
		if err != nil {
			return err
		}

		if resp.GetStatus() == pb.JoinRoomResponse_MATCHED {
			c.room = build.Room(resp.GetRoom())
			c.me = build.User(resp.GetMe())
			fmt.Printf("Matched room_id=%d\n", resp.GetRoom().GetId())
			return nil
		} else if resp.GetStatus() == pb.JoinRoomResponse_WAITING {
			fmt.Println("Waiting matching...")
		}
	}
}

func (c *Chat) chat(ctx context.Context, cli pb.ChatServiceClient) error {
	ct, cancel := context.WithCancel(ctx)
	defer cancel()

	stream, err := cli.Chat(ct)
	if err != nil {
		return err
	}
	defer stream.CloseSend()

	go func() {
		if err := c.send(ct, stream); err != nil {
			cancel()
		}
	}()

	if err = c.receive(ct, stream); err != nil {
		cancel()
		return err
	}

	return nil
}

func (c *Chat) send(ctx context.Context, stream pb.ChatService_ChatClient) error {
	for {
		c.RLock()

		if c.finished {
			c.RUnlock()
			return nil

		} else if !c.started {
			err := stream.Send(&pb.ChatRequest{
				RoomId: c.room.ID,
				User: &pb.User{
					Id:       c.me.ID,
					NickName: c.me.NickName,
				},
				Action: &pb.ChatRequest_Start{
					Start: &pb.ChatRequest_StartAction{},
				},
			})

			c.RUnlock()
			if err != nil {
				return err
			}

			for {
				c.RLock()
				if c.started {
					c.RUnlock()
					fmt.Println("chat start!!!")
					break
				}
				c.RUnlock()
				fmt.Println("Waiting until another user ready")
				time.Sleep(1 * time.Second)
			}
		} else {

			c.RUnlock()
			fmt.Print("Inpurt Your Message!")
			stdin := bufio.NewScanner(os.Stdin)
			stdin.Scan()

			text := stdin.Text()

			go func() {
				err := stream.Send(&pb.ChatRequest{
					RoomId: c.me.ID,
					User: &pb.User{
						Id:       c.me.ID,
						NickName: c.me.NickName,
					},
					Action: &pb.ChatRequest_Talk{
						Talk: &pb.ChatRequest_TalkAction{
							Message: text,
						},
					},
				})
				if err != nil {
					fmt.Println(err)
				}
			}()

			ch := make(chan int)
			go func(ch chan int) {
				fmt.Println("")
				for i := 0; i < 5; i++ {
					fmt.Printf("freezing in %d second.\n", (5 - i))
					time.Sleep(1 * time.Second)
				}
				fmt.Println("")
				ch <- 0
			}(ch)
			<-ch
		}

		select {
		case <-ctx.Done():
			return nil
		default:
		}
	}
}

func (c *Chat) receive(ctx context.Context, stream pb.ChatService_ChatClient) error {
	for {
		res, err := stream.Recv()
		if err != nil {
			return nil
		}

		c.Lock()
		switch res.GetEvent().(type) {
		case *pb.ChatResponse_Waiting:
			// 待機中
		case *pb.ChatResponse_Ready:
			// 開始
			c.started = true
		case *pb.ChatResponse_Chated:
			chatLogs := res.GetChated().ChatLogs
			fmt.Println("chatLogs", chatLogs)
		case *pb.ChatResponse_Finished:
			c.finished = true
			c.Unlock()
			return nil
		}
		c.Unlock()

		select {
		case <- ctx.Done():
			return nil
		default:
		}
	}
}
