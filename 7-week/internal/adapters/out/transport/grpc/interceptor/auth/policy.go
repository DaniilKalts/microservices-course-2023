package auth

import (
	"fmt"
	"strings"

	authv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/7-week/gen/grpc/user/v1"
)

type AccessGroup string

const (
	AccessGroupPublic        AccessGroup = "public"
	AccessGroupAuthenticated AccessGroup = "authenticated"
	AccessGroupAdmin         AccessGroup = "admin"
)

type MethodGroup struct {
	Group   AccessGroup
	Methods []string
}

type AccessPolicy struct {
	groupByMethod map[string]AccessGroup
}

func NewAccessPolicy(groups ...MethodGroup) (AccessPolicy, error) {
	if len(groups) == 0 {
		return AccessPolicy{}, fmt.Errorf("create access policy: no method groups provided")
	}

	groupByMethod := make(map[string]AccessGroup)

	for _, group := range groups {
		if !isSupportedAccessGroup(group.Group) {
			return AccessPolicy{}, fmt.Errorf("create access policy: unsupported access group %q", group.Group)
		}

		if len(group.Methods) == 0 {
			return AccessPolicy{}, fmt.Errorf("create access policy: group %q has no methods", group.Group)
		}

		for _, method := range group.Methods {
			method = strings.TrimSpace(method)
			if method == "" {
				return AccessPolicy{}, fmt.Errorf("create access policy: empty method in group %q", group.Group)
			}

			if existingGroup, exists := groupByMethod[method]; exists {
				return AccessPolicy{}, fmt.Errorf(
					"create access policy: method %q already bound to group %q",
					method,
					existingGroup,
				)
			}

			groupByMethod[method] = group.Group
		}
	}

	return AccessPolicy{groupByMethod: groupByMethod}, nil
}

func DefaultAccessPolicy() (AccessPolicy, error) {
	return NewAccessPolicy(
		MethodGroup{
			Group: AccessGroupPublic,
			Methods: []string{
				authv1.AuthV1_Register_FullMethodName,
				authv1.AuthV1_Login_FullMethodName,
				authv1.AuthV1_Refresh_FullMethodName,
				authv1.AuthV1_Logout_FullMethodName,
				userv1.UserV1_List_FullMethodName,
				userv1.UserV1_Get_FullMethodName,
			},
		},
		MethodGroup{
			Group: AccessGroupAdmin,
			Methods: []string{
				userv1.UserV1_Create_FullMethodName,
				userv1.UserV1_Update_FullMethodName,
				userv1.UserV1_Delete_FullMethodName,
			},
		},
		MethodGroup{
			Group: AccessGroupAuthenticated,
			Methods: []string{
				userv1.ProfileV1_GetProfile_FullMethodName,
				userv1.ProfileV1_UpdateProfile_FullMethodName,
				userv1.ProfileV1_DeleteProfile_FullMethodName,
			},
		},
	)
}

func (p AccessPolicy) GroupForMethod(fullMethod string) (AccessGroup, bool) {
	group, ok := p.groupByMethod[fullMethod]
	return group, ok
}

func (p AccessPolicy) IsEmpty() bool {
	return len(p.groupByMethod) == 0
}

func isSupportedAccessGroup(group AccessGroup) bool {
	switch group {
	case AccessGroupPublic, AccessGroupAuthenticated, AccessGroupAdmin:
		return true
	default:
		return false
	}
}
