package build

import (
	"chat/chat"
	"chat/gen/pb"
)

func Room(r *pb.Room) *chat.Room {
	return &chat.Room{
		ID:    r.Id,
		Host:  User(r.GetHost()),
		Guest: User(r.GetGuest()),
	}
}

func User(u *pb.User) *chat.User {
	return &chat.User{
		ID:       u.GetId(),
		NickName: u.GetNickName(),
	}
}
