package dao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/HelliWrold1/cloud/internal/cache"
	"github.com/HelliWrold1/cloud/internal/model"

	cacheBase "github.com/zhufuyi/sponge/pkg/cache"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/spf13/cast"
	"golang.org/x/sync/singleflight"
	"gorm.io/gorm"
)

var _ UserDao = (*userDao)(nil)

// UserDao defining the dao interface
type UserDao interface {
	Create(ctx context.Context, table *model.User) error
	DeleteByID(ctx context.Context, id uint64) error
	DeleteByIDs(ctx context.Context, ids []uint64) error
	UpdateByID(ctx context.Context, table *model.User) error
	GetByID(ctx context.Context, id uint64) (*model.User, error)
	GetByIDs(ctx context.Context, ids []uint64) ([]*model.User, error)
	GetByColumns(ctx context.Context, params *query.Params) ([]*model.User, int64, error)
	ExistUserByUsername(ctx context.Context, username string) (*model.User, bool)
	UpdateByIDPasswordToNew(ctx context.Context, uid uint64, newPwd string) error
	QueryUserByUsername(ctx context.Context, username string) (*model.User, error)
}

type userDao struct {
	db    *gorm.DB
	cache cache.UserCache
	sfg   *singleflight.Group
}

// NewUserDao creating the dao interface
func NewUserDao(db *gorm.DB, cache cache.UserCache) UserDao {
	return &userDao{db: db, cache: cache, sfg: new(singleflight.Group)}
}

// Create a record, insert the record and the id value is written back to the table
func (d *userDao) Create(ctx context.Context, table *model.User) error {
	err := d.db.WithContext(ctx).Create(table).Error
	_ = d.cache.Del(ctx, table.ID)
	return err
}

// DeleteByID delete a record based on id
func (d *userDao) DeleteByID(ctx context.Context, id uint64) error {
	err := d.db.WithContext(ctx).Where("id = ?", id).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.cache.Del(ctx, id)

	return nil
}

// DeleteByIDs batch delete multiple records
func (d *userDao) DeleteByIDs(ctx context.Context, ids []uint64) error {
	err := d.db.WithContext(ctx).Where("id IN (?)", ids).Delete(&model.User{}).Error
	if err != nil {
		return err
	}

	// delete cache
	for _, id := range ids {
		_ = d.cache.Del(ctx, id)
	}

	return nil
}

// UpdateByID update records by id
func (d *userDao) UpdateByID(ctx context.Context, table *model.User) error {
	if table.ID < 1 {
		return errors.New("id cannot be 0")
	}

	update := map[string]interface{}{}

	if table.Username != "" {
		update["user"] = table.Username
	}
	if table.Password != "" {
		update["password"] = table.Password
	}

	err := d.db.WithContext(ctx).Model(table).Updates(update).Error
	if err != nil {
		return err
	}

	// delete cache
	_ = d.cache.Del(ctx, table.ID)

	return nil
}

// GetByID get a record based on id
func (d *userDao) GetByID(ctx context.Context, id uint64) (*model.User, error) {
	record, err := d.cache.Get(ctx, id)
	if err == nil {
		return record, nil
	}

	if errors.Is(err, model.ErrCacheNotFound) {
		// for the same id, prevent high concurrent simultaneous access to mysql
		val, err, _ := d.sfg.Do(utils.Uint64ToStr(id), func() (interface{}, error) { //nolint
			table := &model.User{}
			err = d.db.WithContext(ctx).Where("id = ?", id).First(table).Error
			if err != nil {
				// if data is empty, set not found cache to prevent cache penetration, default expiration time 10 minutes
				if errors.Is(err, model.ErrRecordNotFound) {
					err = d.cache.SetCacheWithNotFound(ctx, id)
					if err != nil {
						return nil, err
					}
					return nil, model.ErrRecordNotFound
				}
				return nil, err
			}
			// set cache
			err = d.cache.Set(ctx, id, table, cacheBase.DefaultExpireTime)
			if err != nil {
				return nil, fmt.Errorf("cache.Set error: %v, id=%d", err, id)
			}
			return table, nil
		})
		if err != nil {
			return nil, err
		}
		table, ok := val.(*model.User)
		if !ok {
			return nil, model.ErrRecordNotFound
		}
		return table, nil
	} else if errors.Is(err, cacheBase.ErrPlaceholder) {
		return nil, model.ErrRecordNotFound
	}

	// fail fast, if cache error return, don't request to db
	return nil, err
}

// GetByIDs get multiple rows by ids
func (d *userDao) GetByIDs(ctx context.Context, ids []uint64) ([]*model.User, error) {
	records := []*model.User{}

	itemMap, err := d.cache.MultiGet(ctx, ids)
	if err != nil {
		return nil, err
	}

	var missedIDs []uint64
	for _, id := range ids {
		item, ok := itemMap[cast.ToString(id)]
		if !ok {
			missedIDs = append(missedIDs, id)
			continue
		}
		records = append(records, item)
	}

	// get missed data
	if len(missedIDs) > 0 {
		// find the id of an active placeholder, i.e. an id that does not exist in mysql
		var realMissedIDs []uint64
		for _, id := range missedIDs {
			_, err = d.cache.Get(ctx, id)
			if errors.Is(err, cacheBase.ErrPlaceholder) {
				continue
			} else {
				realMissedIDs = append(realMissedIDs, id)
			}
		}

		if len(realMissedIDs) > 0 {
			var missedData []*model.User
			err = d.db.WithContext(ctx).Where("id IN (?)", realMissedIDs).Find(&missedData).Error
			if err != nil {
				return nil, err
			}

			if len(missedData) > 0 {
				records = append(records, missedData...)
				err = d.cache.MultiSet(ctx, missedData, time.Hour*24)
				if err != nil {
					return nil, err
				}
			} else {
				for _, id := range realMissedIDs {
					_ = d.cache.SetCacheWithNotFound(ctx, id)
				}
			}
		}
	}
	return records, nil
}

// GetByColumns filter multiple rows based on paging and column information
//
// params includes paging parameters and query parameters
// paging parameters (required):
//
//	page: page number, starting from 0
//	size: lines per page
//	sort: sort fields, default is id backwards, you can add - sign before the field to indicate reverse order, no - sign to indicate ascending order, multiple fields separated by comma
//
// query parameters (not required):
//
//	name: column name
//	exp: expressions, which default to = when the value is null, have =, ! =, >, >=, <, <=, like
//	value: column name
//	logic: logical type, defaults to and when value is null, only &(and), ||(or)
//
// example: search for a male over 20 years of age
//
//	params = &query.Params{
//	    Page: 0,
//	    Size: 20,
//	    Columns: []query.Column{
//		{
//			Name:    "age",
//			Exp: ">",
//			Value:   20,
//		},
//		{
//			Name:  "gender",
//			Value: "male",
//		},
//	}
func (d *userDao) GetByColumns(ctx context.Context, params *query.Params) ([]*model.User, int64, error) {
	queryStr, args, err := params.ConvertToGormConditions()
	if err != nil {
		return nil, 0, errors.New("query params error: " + err.Error())
	}

	var total int64
	if params.Sort != "ignore count" { // determine if count is required
		err = d.db.WithContext(ctx).Model(&model.User{}).Select([]string{"id"}).Where(queryStr, args...).Count(&total).Error
		if err != nil {
			return nil, 0, err
		}
		if total == 0 {
			return nil, total, nil
		}
	}

	records := []*model.User{}
	order, limit, offset := params.ConvertToPage()
	err = d.db.WithContext(ctx).Order(order).Limit(limit).Offset(offset).Where(queryStr, args...).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	return records, total, err
}

func (d *userDao) ExistUserByUsername(ctx context.Context, username string) (*model.User, bool) {
	record := &model.User{}

	err := d.db.WithContext(ctx).Where("user = ?", username).Find(record).Error
	if err != nil {
		return nil, false
	}
	return record, true
}

func (d *userDao) UpdateByIDPasswordToNew(ctx context.Context, uid uint64, newPwd string) error {
	err := d.db.WithContext(ctx).Update("password", newPwd).Where("id = ?", uid).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *userDao) QueryUserByUsername(ctx context.Context, username string) (*model.User, error) {
	record := &model.User{}
	err := d.db.WithContext(ctx).First(record).Where("username = ?", username).Error
	if err != nil {
		return nil, err
	}
	return record, nil
}
