package client

import (
	"chat/build"
	"chat/chat"
	"chat/gen/pb"
	"context"
	"fmt"
	"sync"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Chat struct {
	sync.RWMutex
	room *chat.Room
	me   *chat.User
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

	return nil
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
