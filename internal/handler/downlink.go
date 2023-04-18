package handler

import (
	"errors"

	"github.com/HelliWrold1/cloud/internal/cache"
	"github.com/HelliWrold1/cloud/internal/dao"
	"github.com/HelliWrold1/cloud/internal/ecode"
	"github.com/HelliWrold1/cloud/internal/model"
	"github.com/HelliWrold1/cloud/internal/types"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/mysql/query"
	"github.com/zhufuyi/sponge/pkg/utils"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
)

var _ DownlinkHandler = (*downlinkHandler)(nil)

// DownlinkHandler defining the handler interface
type DownlinkHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	ListByIDs(c *gin.Context)
	List(c *gin.Context)
}

type downlinkHandler struct {
	iDao dao.DownlinkDao
}

// NewDownlinkHandler creating the handler interface
func NewDownlinkHandler() DownlinkHandler {
	return &downlinkHandler{
		iDao: dao.NewDownlinkDao(
			model.GetDB(),
			cache.NewDownlinkCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create downlink
// @Description submit information to create downlink
// @Tags downlink
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.CreateDownlinkRequest true "downlink information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlink [post]
func (h *downlinkHandler) Create(c *gin.Context) {
	form := &types.CreateDownlinkRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	downlink := &model.Downlink{}
	err = copier.Copy(downlink, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateDownlink)
		return
	}

	err = h.iDao.Create(c.Request.Context(), downlink)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": downlink.ID})
}

// DeleteByID delete a record by ID
// @Summary delete downlink
// @Description delete downlink by id
// @Tags downlink
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlink/{id} [delete]
func (h *downlinkHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getDownlinkIDFromPath(c)
	if isAbort {
		return
	}

	err := h.iDao.DeleteByID(c.Request.Context(), id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// DeleteByIDs delete records by multiple id
// @Summary delete downlinks by multiple id
// @Description delete downlinks by multiple id using a post request
// @Tags downlink
// @Security BearerTokenAuth
// @Param data body types.DeleteDownlinksByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlinks/delete/ids [post]
func (h *downlinkHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteDownlinksByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	err = h.iDao.DeleteByIDs(c.Request.Context(), form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update downlink information
// @Description update downlink information by id
// @Tags downlink
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Param data body types.UpdateDownlinkByIDRequest true "downlink information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlink/{id} [put]
func (h *downlinkHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getDownlinkIDFromPath(c)
	if isAbort {
		return
	}

	form := &types.UpdateDownlinkByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	downlink := &model.Downlink{}
	err = copier.Copy(downlink, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateDownlink)
		return
	}

	err = h.iDao.UpdateByID(c.Request.Context(), downlink)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get downlink details
// @Description get downlink details by id
// @Tags downlink
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlink/{id} [get]
func (h *downlinkHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getDownlinkIDFromPath(c)
	if isAbort {
		return
	}

	downlink, err := h.iDao.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.GetDownlinkByIDRespond{}
	err = copier.Copy(data, downlink)
	if err != nil {
		response.Error(c, ecode.ErrGetDownlink)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"downlink": data})
}

// ListByIDs get records by multiple id
// @Summary get downlinks by multiple id
// @Description get downlinks by multiple id using a post request
// @Tags downlink
// @Security BearerTokenAuth
// @Param data body types.GetDownlinksByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlinks/ids [post]
func (h *downlinkHandler) ListByIDs(c *gin.Context) {
	form := &types.GetDownlinksByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	downlinks, err := h.iDao.GetByIDs(c.Request.Context(), form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())

		return
	}

	data, err := convertDownlinks(downlinks)
	if err != nil {
		response.Error(c, ecode.ErrListDownlink)
		return
	}

	response.Success(c, gin.H{
		"downlinks": data,
	})
}

// List Get multiple records by query parameters
// @Summary get a list of downlinks
// @Description paging and conditional fetching of downlinks lists using post requests
// @Tags downlink
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.Result{}
// @Router /api/v1/downlinks [post]
func (h *downlinkHandler) List(c *gin.Context) {
	form := &types.GetDownlinksRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	downlinks, total, err := h.iDao.GetByColumns(c.Request.Context(), &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertDownlinks(downlinks)
	if err != nil {
		response.Error(c, ecode.ErrListDownlink)
		return
	}

	response.Success(c, gin.H{
		"downlinks": data,
		"total":     total,
	})
}

func getDownlinkIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return "", 0, true
	}

	return idStr, id, false
}

func convertDownlinks(fromValues []*model.Downlink) ([]*types.GetDownlinkByIDRespond, error) {
	toValues := []*types.GetDownlinkByIDRespond{}
	for _, v := range fromValues {
		data := &types.GetDownlinkByIDRespond{}
		err := copier.Copy(data, v)
		if err != nil {
			return nil, err
		}
		data.ID = utils.Uint64ToStr(v.ID)
		toValues = append(toValues, data)
	}

	return toValues, nil
}
