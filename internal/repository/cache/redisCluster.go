package cache

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	//"github.com/go-rbac/internal/core"

	redis "github.com/redis/go-redis/v9"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var childLogger = log.With().Str("repository/cache", "Redis").Logger()

type CacheService struct {
	cache *redis.ClusterClient
}

func NewClusterCache(ctx context.Context, options *redis.ClusterOptions) *CacheService {
	childLogger.Debug().Msg("NewClusterCache")
	childLogger.Debug().Interface("option.Addrs: ", options.Addrs).Msg("")

	redisClient := redis.NewClusterClient(options)
	return &CacheService{
		cache: redisClient,
	}
}

func (s *CacheService) Ping(ctx context.Context) (string, error) {
	childLogger.Debug().Msg("Ping")

	status, err := s.cache.Ping(ctx).Result()
	if err != nil {
		return "", err
	}
	return status, nil
}

func (s *CacheService) Get(ctx context.Context, key string) (interface{}, error) {
	childLogger.Debug().Msg("Get")

	_, root := xray.BeginSubsegment(ctx, "REDIS.Get")
	defer func() {
		root.Close(nil)
	}()

	res, err := s.cache.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s *CacheService) Put(ctx context.Context, key string, value interface{}) error {
	childLogger.Debug().Msg("Put")
	//childLogger.Debug().Str("====> key : ",key).Interface("| Put : ",value).Msg("")

	_, root := xray.BeginSubsegment(ctx, "REDIS.Put")
	defer func() {
		root.Close(nil)
	}()

	value_json, err := json.Marshal(value)
    if err != nil {
       return err
    }

	status := s.cache.Set(ctx, key, value_json, 0)

	return status.Err()
}
