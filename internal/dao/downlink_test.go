package dao

import (
	"context"
	"testing"
	"time"

	"github.com/HelliWrold1/cloud/internal/cache"
	"github.com/HelliWrold1/cloud/internal/model"

	"github.com/zhufuyi/sponge/pkg/gotest"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func newDownlinkDao() *gotest.Dao {
	testData := &model.Downlink{}
	testData.ID = 1
	testData.CreatedAt = time.Now()
	testData.UpdatedAt = testData.CreatedAt

	// init mock cache
	//c := gotest.NewCache(map[string]interface{}{"no cache": testData}) // to test mysql, disable caching
	c := gotest.NewCache(map[string]interface{}{utils.Uint64ToStr(testData.ID): testData})
	c.ICache = cache.NewDownlinkCache(&model.CacheType{
		CType: "redis",
		Rdb:   c.RedisClient,
	})

	// init mock dao
	d := gotest.NewDao(c, testData)
	d.IDao = NewDownlinkDao(d.DB, c.ICache.(cache.DownlinkCache))

	return d
}

func Test_downlinkDao_Create(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("INSERT INTO .*").
		WithArgs(d.GetAnyArgs(testData)...).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(DownlinkDao).Create(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_downlinkDao_DeleteByID(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	testData.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: false,
	}

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(DownlinkDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(DownlinkDao).DeleteByID(d.Ctx, 0)
	assert.Error(t, err)
}

func Test_downlinkDao_DeleteByIDs(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	testData.DeletedAt = gorm.DeletedAt{
		Time:  time.Now(),
		Valid: false,
	}

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(int64(testData.ID), 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(DownlinkDao).DeleteByID(d.Ctx, testData.ID)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(DownlinkDao).DeleteByIDs(d.Ctx, []uint64{0})
	assert.Error(t, err)
}

func Test_downlinkDao_UpdateByID(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	d.SQLMock.ExpectBegin()
	d.SQLMock.ExpectExec("UPDATE .*").
		WithArgs(d.AnyTime, testData.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	d.SQLMock.ExpectCommit()

	err := d.IDao.(DownlinkDao).UpdateByID(d.Ctx, testData)
	if err != nil {
		t.Fatal(err)
	}

	// zero id error
	err = d.IDao.(DownlinkDao).UpdateByID(d.Ctx, &model.Downlink{})
	assert.Error(t, err)

}

func Test_downlinkDao_GetByID(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(DownlinkDao).GetByID(d.Ctx, testData.ID) // notfound
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// notfound error
	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(2).
		WillReturnRows(rows)
	_, err = d.IDao.(DownlinkDao).GetByID(d.Ctx, 2)
	assert.Error(t, err)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(3, 4).
		WillReturnRows(rows)
	_, err = d.IDao.(DownlinkDao).GetByID(d.Ctx, 4)
	assert.Error(t, err)
}

func Test_downlinkDao_GetByIDs(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").
		WithArgs(testData.ID).
		WillReturnRows(rows)

	_, err := d.IDao.(DownlinkDao).GetByIDs(d.Ctx, []uint64{testData.ID})
	if err != nil {
		t.Fatal(err)
	}

	_, err = d.IDao.(DownlinkDao).GetByIDs(d.Ctx, []uint64{111})
	assert.Error(t, err)

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}
}

func Test_downlinkDao_GetByColumns(t *testing.T) {
	d := newDownlinkDao()
	defer d.Close()
	testData := d.TestData.(*model.Downlink)

	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at"}).
		AddRow(testData.ID, testData.CreatedAt, testData.UpdatedAt)

	d.SQLMock.ExpectQuery("SELECT .*").WillReturnRows(rows)

	_, _, err := d.IDao.(DownlinkDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Sort: "ignore count", // ignore test count(*)
	})
	if err != nil {
		t.Fatal(err)
	}

	err = d.SQLMock.ExpectationsWereMet()
	if err != nil {
		t.Fatal(err)
	}

	// err test
	_, _, err = d.IDao.(DownlinkDao).GetByColumns(d.Ctx, &query.Params{
		Page: 0,
		Size: 10,
		Columns: []query.Column{
			{
				Name:  "id",
				Exp:   "<",
				Value: 0,
			},
		},
	})
	assert.Error(t, err)

	// error test
	dao := &downlinkDao{}
	_, _, err = dao.GetByColumns(context.Background(), &query.Params{Columns: []query.Column{{}}})
	t.Log(err)
}
