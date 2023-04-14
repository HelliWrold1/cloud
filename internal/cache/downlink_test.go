package cache

import (
	"testing"
	"time"

	"github.com/HelliWrold1/cloud/internal/model"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/stretchr/testify/assert"
)

func newDownlinkCache() *gotest.Cache {
	record1 := &model.Downlink{}
	record1.ID = 1
	record2 := &model.Downlink{}
	record2.ID = 2
	testData := map[string]interface{}{
		utils.Uint64ToStr(record1.ID): record1,
		utils.Uint64ToStr(record2.ID): record2,
	}

	c := gotest.NewCache(testData)
	c.ICache = NewDownlinkCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})
	return c
}

func Test_downlinkCache_Set(t *testing.T) {
	c := newDownlinkCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Downlink)
	err := c.ICache.(DownlinkCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	// nil data
	err = c.ICache.(DownlinkCache).Set(c.Ctx, 0, nil, time.Hour)
	assert.NoError(t, err)
}

func Test_downlinkCache_Get(t *testing.T) {
	c := newDownlinkCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Downlink)
	err := c.ICache.(DownlinkCache).Set(c.Ctx, record.ID, record, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DownlinkCache).Get(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, record, got)

	// zero key error
	_, err = c.ICache.(DownlinkCache).Get(c.Ctx, 0)
	assert.Error(t, err)
}

func Test_downlinkCache_MultiGet(t *testing.T) {
	c := newDownlinkCache()
	defer c.Close()

	var testData []*model.Downlink
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Downlink))
	}

	err := c.ICache.(DownlinkCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}

	got, err := c.ICache.(DownlinkCache).MultiGet(c.Ctx, c.GetIDs())
	if err != nil {
		t.Fatal(err)
	}

	expected := c.GetTestData()
	for k, v := range expected {
		assert.Equal(t, got[k], v.(*model.Downlink))
	}
}

func Test_downlinkCache_MultiSet(t *testing.T) {
	c := newDownlinkCache()
	defer c.Close()

	var testData []*model.Downlink
	for _, data := range c.TestDataSlice {
		testData = append(testData, data.(*model.Downlink))
	}

	err := c.ICache.(DownlinkCache).MultiSet(c.Ctx, testData, time.Hour)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_downlinkCache_Del(t *testing.T) {
	c := newDownlinkCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Downlink)
	err := c.ICache.(DownlinkCache).Del(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_downlinkCache_SetCacheWithNotFound(t *testing.T) {
	c := newDownlinkCache()
	defer c.Close()

	record := c.TestDataSlice[0].(*model.Downlink)
	err := c.ICache.(DownlinkCache).SetCacheWithNotFound(c.Ctx, record.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNewDownlinkCache(t *testing.T) {
	c := NewDownlinkCache(&model.CacheType{
		CType: "memory",
	})

	assert.NotNil(t, c)
}
