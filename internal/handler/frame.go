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

var _ FrameHandler = (*frameHandler)(nil)

// FrameHandler defining the handler interface
type FrameHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	ListByIDs(c *gin.Context)
	List(c *gin.Context)
}

type frameHandler struct {
	iDao dao.FrameDao
}

// NewFrameHandler creating the handler interface
func NewFrameHandler() FrameHandler {
	return &frameHandler{
		iDao: dao.NewFrameDao(
			model.GetDB(),
			cache.NewFrameCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create frame
// @Description submit information to create frame
// @Tags frame
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.CreateFrameRequest true "frame information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/frame [post]
func (h *frameHandler) Create(c *gin.Context) {
	form := &types.CreateFrameRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	frame := &model.Frame{}
	err = copier.Copy(frame, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateFrame)
		return
	}

	err = h.iDao.Create(c.Request.Context(), frame)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": frame.ID})
}

// DeleteByID delete a record by ID
// @Summary delete frame
// @Description delete frame by id
// @Tags frame
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Success 200 {object} types.Result{}
// @Router /api/v1/frame/{id} [delete]
func (h *frameHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getFrameIDFromPath(c)
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
// @Summary delete frames by multiple id
// @Description delete frames by multiple id using a post request
// @Tags frame
// @Security BearerTokenAuth
// @Param data body types.DeleteFramesByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/frames/delete/ids [post]
func (h *frameHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteFramesByIDsRequest{}
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
// @Summary update frame information
// @Description update frame information by id
// @Tags frame
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Param data body types.UpdateFrameByIDRequest true "frame information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/frame/{id} [put]
func (h *frameHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getFrameIDFromPath(c)
	if isAbort {
		return
	}

	form := &types.UpdateFrameByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	frame := &model.Frame{}
	err = copier.Copy(frame, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateFrame)
		return
	}

	err = h.iDao.UpdateByID(c.Request.Context(), frame)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get frame details
// @Description get frame details by id
// @Tags frame
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/frame/{id} [get]
func (h *frameHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getFrameIDFromPath(c)
	if isAbort {
		return
	}

	frame, err := h.iDao.GetByID(c.Request.Context(), id)
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

	data := &types.GetFrameByIDRespond{}
	err = copier.Copy(data, frame)
	if err != nil {
		response.Error(c, ecode.ErrGetFrame)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"frame": data})
}

// ListByIDs get records by multiple id
// @Summary get frames by multiple id
// @Description get frames by multiple id using a post request
// @Tags frame
// @Security BearerTokenAuth
// @Param data body types.GetFramesByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/frames/ids [post]
func (h *frameHandler) ListByIDs(c *gin.Context) {
	form := &types.GetFramesByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	frames, err := h.iDao.GetByIDs(c.Request.Context(), form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())

		return
	}

	data, err := convertFrames(frames)
	if err != nil {
		response.Error(c, ecode.ErrListFrame)
		return
	}

	response.Success(c, gin.H{
		"frames": data,
	})
}

// List Get multiple records by query parameters
// @Summary get a list of frames
// @Description paging and conditional fetching of frames lists using post requests
// @Tags frame
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.Result{}
// @Router /api/v1/frames [post]
func (h *frameHandler) List(c *gin.Context) {
	form := &types.GetFramesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	frames, total, err := h.iDao.GetByColumns(c.Request.Context(), &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertFrames(frames)
	if err != nil {
		response.Error(c, ecode.ErrListFrame)
		return
	}

	response.Success(c, gin.H{
		"frames": data,
		"total":  total,
	})
}

func getFrameIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return "", 0, true
	}

	return idStr, id, false
}

func convertFrames(fromValues []*model.Frame) ([]*types.GetFrameByIDRespond, error) {
	toValues := []*types.GetFrameByIDRespond{}
	for _, v := range fromValues {
		data := &types.GetFrameByIDRespond{}
		err := copier.Copy(data, v)
		if err != nil {
			return nil, err
		}
		data.ID = utils.Uint64ToStr(v.ID)
		toValues = append(toValues, data)
	}

	return toValues, nil
}
