package interceptor

import (
	"fmt"
	"strings"
)

type AccessGroup struct {
	Name     string
	RoleIDs  []int32
	IsPublic bool
}

func PublicGroup() AccessGroup {
	return AccessGroup{Name: "public", IsPublic: true}
}

func RoleGroup(name string, roleIDs ...int32) AccessGroup {
	return AccessGroup{Name: name, RoleIDs: roleIDs}
}

func (g AccessGroup) AllowsRole(roleID int32) bool {
	if g.IsPublic {
		return true
	}

	for _, allowed := range g.RoleIDs {
		if allowed == roleID {
			return true
		}
	}

	return false
}

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
		if group.Group.Name == "" {
			return AccessPolicy{}, fmt.Errorf("create access policy: group name is empty")
		}

		if len(group.Methods) == 0 {
			return AccessPolicy{}, fmt.Errorf("create access policy: group %q has no methods", group.Group.Name)
		}

		for _, method := range group.Methods {
			method = strings.TrimSpace(method)
			if method == "" {
				return AccessPolicy{}, fmt.Errorf("create access policy: empty method in group %q", group.Group.Name)
			}

			if existingGroup, exists := groupByMethod[method]; exists {
				return AccessPolicy{}, fmt.Errorf(
					"create access policy: method %q already bound to group %q",
					method,
					existingGroup.Name,
				)
			}

			groupByMethod[method] = group.Group
		}
	}

	return AccessPolicy{groupByMethod: groupByMethod}, nil
}

func (p AccessPolicy) GroupForMethod(fullMethod string) (AccessGroup, bool) {
	group, ok := p.groupByMethod[fullMethod]
	return group, ok
}

func (p AccessPolicy) IsEmpty() bool {
	return len(p.groupByMethod) == 0
}
