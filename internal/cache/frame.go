package cache

import (
	"context"
	"strings"
	"time"

	"github.com/HelliWrold1/cloud/internal/model"

	"github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/encoding"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/spf13/cast"
)

const (
	// PrefixFrameCacheKey cache prefix
	PrefixFrameCacheKey = "frame:"
)

var _ FrameCache = (*frameCache)(nil)

// FrameCache cache interface
type FrameCache interface {
	Set(ctx context.Context, id uint64, data *model.Frame, duration time.Duration) error
	Get(ctx context.Context, id uint64) (ret *model.Frame, err error)
	MultiGet(ctx context.Context, ids []uint64) (map[string]*model.Frame, error)
	MultiSet(ctx context.Context, data []*model.Frame, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// frameCache define a cache struct
type frameCache struct {
	cache cache.Cache
}

// NewFrameCache new a cache
func NewFrameCache(cacheType *model.CacheType) FrameCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	var c cache.Cache
	if strings.ToLower(cacheType.CType) == "redis" {
		c = cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Frame{}
		})
	} else {
		c = cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Frame{}
		})
	}

	return &frameCache{
		cache: c,
	}
}

// GetFrameCacheKey cache key
func (c *frameCache) GetFrameCacheKey(id uint64) string {
	return PrefixFrameCacheKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *frameCache) Set(ctx context.Context, id uint64, data *model.Frame, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetFrameCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *frameCache) Get(ctx context.Context, id uint64) (*model.Frame, error) {
	var data *model.Frame
	cacheKey := c.GetFrameCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *frameCache) MultiSet(ctx context.Context, data []*model.Frame, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetFrameCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *frameCache) MultiGet(ctx context.Context, ids []uint64) (map[string]*model.Frame, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetFrameCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Frame)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.Frame)
	for _, v := range ids {
		val, ok := itemMap[c.GetFrameCacheKey(v)]
		if ok {
			retMap[cast.ToString(v)] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *frameCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetFrameCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *frameCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetFrameCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
