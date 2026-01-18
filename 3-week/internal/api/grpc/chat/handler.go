package chat

import (
	chatv1 "github.com/DaniilKalts/microservices-course-2023/3-week/gen/go/chat/v1"
)

type Implementation struct {
	chatv1.UnimplementedChatV1Server
}

func NewImplementation() *Implementation {
	return &Implementation{}
}
