package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	authv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/auth/v1"
	userv1 "github.com/DaniilKalts/microservices-course-2023/8-week/api/gen/go/user/v1"
)

// TestDefaultAccessPolicyCoverage ensures every registered gRPC method has an
// entry in DefaultAccessPolicy. If a new proto method is added without updating
// the policy, this test will fail.
func TestDefaultAccessPolicyCoverage(t *testing.T) {
	policy, err := DefaultAccessPolicy()
	require.NoError(t, err)

	server := grpc.NewServer()

	authv1.RegisterAuthV1Server(server, nil)
	userv1.RegisterUserV1Server(server, nil)
	userv1.RegisterProfileV1Server(server, nil)

	serviceInfo := server.GetServiceInfo()

	var uncovered []string
	for serviceName, info := range serviceInfo {
		for _, method := range info.Methods {
			fullMethod := "/" + serviceName + "/" + method.Name
			if _, ok := policy.GroupForMethod(fullMethod); !ok {
				uncovered = append(uncovered, fullMethod)
			}
		}
	}

	assert.Emptyf(t, uncovered, "gRPC methods missing from DefaultAccessPolicy (update policy_default.go):\n%v", uncovered)
}
