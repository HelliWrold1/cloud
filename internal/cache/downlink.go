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
	// PrefixDownlinkCacheKey cache prefix
	PrefixDownlinkCacheKey = "downlink:"
)

var _ DownlinkCache = (*downlinkCache)(nil)

// DownlinkCache cache interface
type DownlinkCache interface {
	Set(ctx context.Context, id uint64, data *model.Downlink, duration time.Duration) error
	Get(ctx context.Context, id uint64) (ret *model.Downlink, err error)
	MultiGet(ctx context.Context, ids []uint64) (map[string]*model.Downlink, error)
	MultiSet(ctx context.Context, data []*model.Downlink, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetCacheWithNotFound(ctx context.Context, id uint64) error
}

// downlinkCache define a cache struct
type downlinkCache struct {
	cache cache.Cache
}

// NewDownlinkCache new a cache
func NewDownlinkCache(cacheType *model.CacheType) DownlinkCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""
	var c cache.Cache
	if strings.ToLower(cacheType.CType) == "redis" {
		c = cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Downlink{}
		})
	} else {
		c = cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Downlink{}
		})
	}

	return &downlinkCache{
		cache: c,
	}
}

// GetDownlinkCacheKey cache key
func (c *downlinkCache) GetDownlinkCacheKey(id uint64) string {
	return PrefixDownlinkCacheKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *downlinkCache) Set(ctx context.Context, id uint64, data *model.Downlink, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetDownlinkCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *downlinkCache) Get(ctx context.Context, id uint64) (*model.Downlink, error) {
	var data *model.Downlink
	cacheKey := c.GetDownlinkCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *downlinkCache) MultiSet(ctx context.Context, data []*model.Downlink, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetDownlinkCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *downlinkCache) MultiGet(ctx context.Context, ids []uint64) (map[string]*model.Downlink, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetDownlinkCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Downlink)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[string]*model.Downlink)
	for _, v := range ids {
		val, ok := itemMap[c.GetDownlinkCacheKey(v)]
		if ok {
			retMap[cast.ToString(v)] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *downlinkCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetDownlinkCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetCacheWithNotFound set empty cache
func (c *downlinkCache) SetCacheWithNotFound(ctx context.Context, id uint64) error {
	cacheKey := c.GetDownlinkCacheKey(id)
	err := c.cache.SetCacheWithNotFound(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}
