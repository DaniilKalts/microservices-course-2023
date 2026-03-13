package auth

import (
	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
	domainUser "github.com/DaniilKalts/microservices-course-2023/8-week/internal/domain/user"
)

func DefaultAccessPolicy() (AccessPolicy, error) {
	public := PublicGroup()
	authenticated := RoleGroup("authenticated", int32(domainUser.RoleUser), int32(domainUser.RoleAdmin))
	admin := RoleGroup("admin", int32(domainUser.RoleAdmin))

	return NewAccessPolicy(
		AccessRule{
			Group: public,
			Methods: []string{
				authv1.AuthV1_Register_FullMethodName,
				authv1.AuthV1_Login_FullMethodName,
				authv1.AuthV1_Refresh_FullMethodName,
				userv1.UserV1_List_FullMethodName,
				userv1.UserV1_Get_FullMethodName,
				"/grpc.health.v1.Health/Check",
				"/grpc.health.v1.Health/Watch",
			},
		},
		AccessRule{
			Group: admin,
			Methods: []string{
				userv1.UserV1_Create_FullMethodName,
				userv1.UserV1_Update_FullMethodName,
				userv1.UserV1_Delete_FullMethodName,
			},
		},
		AccessRule{
			Group: authenticated,
			Methods: []string{
				authv1.AuthV1_Logout_FullMethodName,
				userv1.ProfileV1_GetProfile_FullMethodName,
				userv1.ProfileV1_UpdateProfile_FullMethodName,
				userv1.ProfileV1_DeleteProfile_FullMethodName,
			},
		},
	)
}
