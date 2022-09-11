package handler

import (
	"chat/build"
	"chat/chat"
	"chat/gen/pb"
	"context"
	"fmt"
	"sync"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Match struct {
	sync.RWMutex
	Rooms     map[int32]*chat.Room
	maxUserID int32
}

func NewMatch() *Match {
	return &Match{
		Rooms: make(map[int32]*chat.Room),
	}
}

func (m *Match) JoinRoom(_ *pb.JoinRoomRequest, stream pb.MatchService_JoinRoomServer) error {
	ctx, cancel := context.WithTimeout(stream.Context(), 2*time.Minute)
	defer cancel()

	m.Lock()

	m.maxUserID++
	me := &chat.User{
		ID: m.maxUserID,
	}

	// 合いているチャットルームを探す
	for _, room := range m.Rooms {
		if room.Guest == nil {
			room.Guest = me
			stream.Send(&pb.JoinRoomResponse{
				Room:   build.PBRoom(room),
				Me:     build.PBUser(me),
				Status: pb.JoinRoomResponse_MATCHED,
			})
			m.Unlock()
			fmt.Printf("matched room_id=%v\n", room.ID)
			return nil
		}
	}

	// 合いている部屋がない場合は部屋を作る
	room := &chat.Room{
		ID:   int32(len(m.Rooms)) + 1,
		Host: me,
	}
	m.Rooms[room.ID] = room
	m.Unlock()

	stream.Send(&pb.JoinRoomResponse{
		Room:   build.PBRoom(room),
		Status: pb.JoinRoomResponse_WAITING,
	})

	ch := make(chan int)
	go func(ch chan<- int) {
		for {
			m.RLock()
			guest := room.Guest
			m.RUnlock()

			if guest != nil {
				stream.Send(&pb.JoinRoomResponse{
					Room:   build.PBRoom(room),
					Me:     build.PBUser(room.Host),
					Status: pb.JoinRoomResponse_MATCHED,
				})
				ch <- 0
				break
			}
			time.Sleep(1 * time.Second)

			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}(ch)

	select {
	case <-ch:
	case <-ctx.Done():
		return status.Errorf(codes.DeadlineExceeded, "マッチングされませんでした")
	}

	return nil
}
