package model

import (
	"sync"
	"time"

	"github.com/HelliWrold1/cloud/internal/config"

	"github.com/zhufuyi/sponge/pkg/goredis"
	"github.com/zhufuyi/sponge/pkg/mysql"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

var (
	// ErrCacheNotFound No hit cache
	ErrCacheNotFound = redis.Nil

	// ErrRecordNotFound no records found
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

var (
	DB    *gorm.DB
	once1 sync.Once

	redisCli *redis.Client
	once2    sync.Once

	cacheType *CacheType
	once3     sync.Once
)

// InitMysql connect mysql
func InitMysql() {
	opts := []mysql.Option{
		mysql.WithSlowThreshold(time.Duration(config.Get().Mysql.SlowThreshold) * time.Millisecond),
		mysql.WithMaxIdleConns(config.Get().Mysql.MaxIdleConns),
		mysql.WithMaxOpenConns(config.Get().Mysql.MaxOpenConns),
		mysql.WithConnMaxLifetime(time.Duration(config.Get().Mysql.ConnMaxLifetime) * time.Minute),
	}
	if config.Get().Mysql.EnableLog {
		opts = append(opts, mysql.WithLog())
	}

	if config.Get().App.EnableTrace {
		opts = append(opts, mysql.WithEnableTrace())
	}

	var err error
	DB, err = mysql.Init(config.Get().Mysql.Dsn, opts...)
	if err != nil {
		panic("mysql.Init error: " + err.Error())
	}
	migration()
}

// GetDB get db
func GetDB() *gorm.DB {
	if DB == nil {
		once1.Do(func() {
			InitMysql()
		})
	}

	return DB
}

// CloseMysql close mysql
func CloseMysql() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	if sqlDB != nil {
		return sqlDB.Close()
	}
	return nil
}

// ------------------------------------------------------------------------------------------

// CacheType cache type
type CacheType struct {
	CType string        // cache type  memory or redis
	Rdb   *redis.Client // if CType=redis, Rdb cannot be empty
}

// InitCache initial cache
func InitCache(cType string) {
	cacheType = &CacheType{
		CType: cType,
	}

	if cType == "redis" {
		cacheType.Rdb = GetRedisCli()
	}
}

// GetCacheType get cacheType
func GetCacheType() *CacheType {
	if cacheType == nil {
		once3.Do(func() {
			InitCache(config.Get().App.CacheType)
		})
	}

	return cacheType
}

// InitRedis connect redis
func InitRedis() {
	opts := []goredis.Option{
		goredis.WithDialTimeout(time.Duration(config.Get().Redis.DialTimeout) * time.Second),
		goredis.WithReadTimeout(time.Duration(config.Get().Redis.ReadTimeout) * time.Second),
		goredis.WithWriteTimeout(time.Duration(config.Get().Redis.WriteTimeout) * time.Second),
	}
	if config.Get().App.EnableTrace {
		opts = append(opts, goredis.WithEnableTrace())
	}

	var err error
	redisCli, err = goredis.Init(config.Get().Redis.Dsn, opts...)
	if err != nil {
		panic("goredis.Init error: " + err.Error())
	}
}

// GetRedisCli get redis client
func GetRedisCli() *redis.Client {
	if redisCli == nil {
		once2.Do(func() {
			InitRedis()
		})
	}

	return redisCli
}

// CloseRedis close redis
func CloseRedis() error {
	if redisCli == nil {
		return nil
	}

	err := redisCli.Close()
	if err != nil && err.Error() != redis.ErrClosed.Error() {
		return err
	}

	return nil
}
