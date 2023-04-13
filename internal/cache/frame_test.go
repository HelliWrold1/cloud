package cache

import (
	"testing"
	"time"

	"github.com/HelliWrold1/cloud/internal/model"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func newFrameCache() *gotest.Cache {
	record1 := &model.Frame{}
	record1.ID = 1
	record2 := &model.Frame{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewFrameCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_frameCache_Set(t *testing.T) {
	c := newFrameCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Frame)
	err := c.ICache.(FrameCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(FrameCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_frameCache_Get(t *testing.T) {
	c := newFrameCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Frame)
	err := c.ICache.(FrameCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(FrameCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(FrameCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_frameCache_MultiGet(t *testing.T) {
	c := newFrameCache()
	defer c.Close()

	var testData []*model.Frame
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Frame))
	}

	err := c.ICache.(FrameCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(FrameCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[k], v.(*model.Frame))
	}
}

func Test_frameCache_MultiSet(t *testing.T) {
	c := newFrameCache()
	defer c.Close()

	var testData []*model.Frame
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Frame))
	}

	err := c.ICache.(FrameCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_frameCache_Del(t *testing.T) {
	c := newFrameCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Frame)
	err := c.ICache.(FrameCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_frameCache_SetCacheWithNotFound(t *testing.T) {
	c := newFrameCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Frame)
	err := c.ICache.(FrameCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewFrameCache(t *testing.T) {
	c := NewFrameCache(&model.CacheType{
		CType: "memory",
	})

	assert.NotNil(t, c)
}
