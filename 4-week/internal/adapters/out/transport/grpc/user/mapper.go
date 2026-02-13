package user

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	userv1 "github.com/DaniilKalts/microservices-course-2023/4-week/gen/go/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/4-week/internal/domain/user"
)

func toDomainFromCreate(req *userv1.CreateRequest) *domainUser.Entity {
	return &domainUser.Entity{
		Name:  req.GetName(),
		Email: req.GetEmail(),
		Role:  domainUser.Role(req.GetRole()),
	}
}

func toDomainPatchFromUpdate(req *userv1.UpdateRequest) *domainUser.UpdatePatch {
	patch := &domainUser.UpdatePatch{}
	if req.GetName() != nil {
		name := req.GetName().GetValue()
		patch.Name = &name
	}
	if req.GetEmail() != nil {
		email := req.GetEmail().GetValue()
		patch.Email = &email
	}

	return patch
}

func toProtoUser(user *domainUser.Entity) *userv1.User {
	protoUser := &userv1.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      userv1.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
	}

	if user.UpdatedAt != nil {
		protoUser.UpdatedAt = timestamppb.New(*user.UpdatedAt)
	}

	return protoUser
}
