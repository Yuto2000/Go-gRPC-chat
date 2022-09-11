package build

import (
	"chat/chat"
	"chat/gen/pb"
)

func PBRoom(r *chat.Room) *pb.Room {
	return &pb.Room{
		Id: r.ID,
		Host: PBUser(r.Host),
		Guest: PBUser(r.Guest),
	}
}

func PBUser(u *chat.User) *pb.User {
	if u == nil {
		return nil
	}
	return &pb.User{
		Id:       u.ID,
		NickName: u.NickName,
	}
}