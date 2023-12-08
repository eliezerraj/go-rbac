package service

import(
	"context"

	"github.com/go-rbac/internal/core"
	"github.com/aws/aws-xray-sdk-go/xray"
)

func (s *RBACService) Authentication(ctx context.Context, policy core.PolicyData) (bool, error) {
	childLogger.Debug().Msg("Authentication")

	_, root := xray.BeginSubsegment(ctx, "RBACService.Authentication")
	defer func() {
		root.Close(nil)
	}()

	// Put to the cache
	key := "policy:" + policy.Policy.Name

	err := s.cacheRedis.Put(ctx, key, policy.Policy)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return false, err
	}
	
	return true, nil
}

func (s *RBACService) Enforce(ctx context.Context, user core.User, resource core.Resource, action string) (bool, error) {
	childLogger.Debug().Msg("Enforce")

	_, root := xray.BeginSubsegment(ctx, "RBACService.Enforce")
	defer func() {
		root.Close(nil)
	}()

	return true, nil
}