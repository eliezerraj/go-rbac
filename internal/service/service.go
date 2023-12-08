package service

import (
		"context"
		"errors"
		"encoding/json"
		"fmt"

		"github.com/mitchellh/mapstructure"
		"github.com/rs/zerolog/log"
		"github.com/go-rbac/internal/repository/cache"
		"github.com/go-rbac/internal/core"

		"github.com/aws/aws-xray-sdk-go/xray"
)

var childLogger = log.With().Str("service", "RBACService").Logger()

type RBACService struct {
	cacheRedis	*cache.CacheService
}

func NewRBACService(cacheRedis	*cache.CacheService) *RBACService{
	childLogger.Debug().Msg("NewRBACService")

	return &RBACService{
		cacheRedis: cacheRedis,
	}
}

func (s *RBACService) PutPolicy(ctx context.Context, policy core.PolicyData) (bool, error) {
	childLogger.Debug().Msg("PutPolicy")

	_, root := xray.BeginSubsegment(ctx, "RBACService.PutPolicy")
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

func (s *RBACService) PutRole(ctx context.Context, role core.RoleData) (bool, error) {
	childLogger.Debug().Msg("PutRole")

	_, root := xray.BeginSubsegment(ctx, "RBACService.PutRole")
	defer func() {
		root.Close(nil)
	}()

	// Put to the cache
	key := "role:" + role.Role.Name

	err := s.cacheRedis.Put(ctx, key, role.Role)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return false, err
	}

	return true, nil
}

func (s *RBACService) GetRole(ctx context.Context, role core.RoleData) (*core.Role, error) {
	childLogger.Debug().Msg("GetRole")

	_, root := xray.BeginSubsegment(ctx, "RBACService.GetRole")
	defer func() {
		root.Close(nil)
	}()

	// Get to the cache
	key := "role:"+ role.Role.Name

	res, err := s.cacheRedis.Get(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return nil, err
	}

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(fmt.Sprint(res)), &jsonMap)

	childLogger.Debug().Interface("@@@@ jsonMap : ", jsonMap).Msg("")

	var role_assert core.Role
	err = mapstructure.Decode(jsonMap, &role_assert)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	return &role_assert, nil
}

func (s *RBACService) GetPolicy(ctx context.Context, policy core.PolicyData) (*core.PolicyData, error) {
	childLogger.Debug().Msg("GetPolicy")

	_, root := xray.BeginSubsegment(ctx, "RBACService.GetPolicy")
	defer func() {
		root.Close(nil)
	}()

	// Get to the cache
	key := "policy:"+ policy.Policy.Name

	res, err := s.cacheRedis.Get(ctx, key)
	if err != nil {
		childLogger.Error().Err(err).Msg(".")
		return nil, err
	}

	//Assertion
	var policy_assert core.PolicyData
	err = mapstructure.Decode(res, &policy_assert)
    if err != nil {
		childLogger.Error().Err(err).Msg("error parse interface")
		return nil, errors.New(err.Error())
    }

	return &policy_assert, nil
}
