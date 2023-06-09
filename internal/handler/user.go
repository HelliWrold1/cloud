package handler

import (
	"errors"
	"github.com/HelliWrold1/cloud/internal/crypt"
	"github.com/zhufuyi/sponge/pkg/jwt"

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

var _ UserHandler = (*userHandler)(nil)

// UserHandler defining the handler interface
type UserHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	ListByIDs(c *gin.Context)
	List(c *gin.Context)
	LoginUser(c *gin.Context)
	UpdateByUserPasswordToNew(c *gin.Context)
	GetUserInfo(c *gin.Context)
	LogoutUser(c *gin.Context)
}

type userHandler struct {
	iDao dao.UserDao
}

// NewUserHandler creating the handler interface
func NewUserHandler() UserHandler {
	return &userHandler{
		iDao: dao.NewUserDao(
			model.GetDB(),
			cache.NewUserCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary register user
// @Description submit information to register user
// @Tags user
// @accept json
// @Produce json
// @Param data body types.CreateUserRequest true "user information"
// @Security BearerTokenAuth
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/register [post]
func (h *userHandler) Create(c *gin.Context) {
	form := &types.CreateUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUser)
		return
	}

	// 查找用户是否已存在
	_, exist := h.iDao.ExistUserByUsername(c.Request.Context(), user.Username)
	if exist {
		logger.Error("User already exist", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrCreateUser, "User already exist")
		return
	}
	// 不存在则创建
	cryptedPwd, err := crypt.GenerateSaltPwd(user.Password) // 加密密码
	if err != nil {
		logger.Error("Crypt error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrCreateUser)
		return
	}
	user.Password = cryptedPwd
	err = h.iDao.Create(c.Request.Context(), user)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrCreateUser)
		return
	}

	token, err := jwt.GenerateToken(utils.Uint64ToStr(user.ID), user.Username)
	if err != nil {
		logger.Error("Token error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrCreateUser)
		return
	}
	response.Success(c, gin.H{"id": user.ID, "username": user.Username, "token": token})
}

// DeleteByID delete a record by ID
// @Summary delete user
// @Description delete user by id
// @Tags user
// @accept json
// @Produce json
// @Param id path string true "id"
// @Security BearerTokenAuth
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/{id} [delete]
func (h *userHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
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
// @Summary delete users by multiple id
// @Description delete users by multiple id using a post request
// @Tags user
// @Param data body types.DeleteUsersByIDsRequest true "id array"
// @Security BearerTokenAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/users/delete/ids [post]
func (h *userHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUsersByIDsRequest{}
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
// @Summary update user information
// @Description update user information by id
// @Tags user
// @accept json
// @Produce json
// @Param id path string true "id"
// @Security BearerTokenAuth
// @Param data body types.UpdateUserByIDRequest true "user information"
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/{id} [put]
func (h *userHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		return
	}

	form := &types.UpdateUserByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateUser)
		return
	}

	user.Password, err = crypt.GenerateSaltPwd(form.Password) // 给新密码加盐
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	err = h.iDao.UpdateByID(c.Request.Context(), user)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get user details
// @Description get user details by id
// @Tags user
// @Security BearerTokenAuth
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Security BearerTokenAuth
// @Router /api/v1/user/{id} [get]
func (h *userHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		return
	}

	user, err := h.iDao.GetByID(c.Request.Context(), id)
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

	data := &types.GetUserByIDRespond{}
	err = copier.Copy(data, user)
	if err != nil {
		response.Error(c, ecode.ErrGetUser)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"user": data})
}

// ListByIDs get records by multiple id
// @Summary get users by multiple id
// @Description get users by multiple id using a post request
// @Tags user
// @Security BearerTokenAuth
// @Param data body types.GetUsersByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Security BearerTokenAuth
// @Router /api/v1/users/ids [post]
func (h *userHandler) ListByIDs(c *gin.Context) {
	form := &types.GetUsersByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	users, err := h.iDao.GetByIDs(c.Request.Context(), form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUsers(users)
	if err != nil {
		response.Error(c, ecode.ErrListUser)
		return
	}

	response.Success(c, gin.H{
		"users": data,
	})
}

// List Get multiple records by query parameters
// @Summary get a list of users
// @Description paging and conditional fetching of users lists using post requests
// @Tags user
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.Result{}
// @Router /api/v1/users [post]
func (h *userHandler) List(c *gin.Context) {
	form := &types.GetUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	users, total, err := h.iDao.GetByColumns(c.Request.Context(), &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUsers(users)
	if err != nil {
		response.Error(c, ecode.ErrListUser)
		return
	}

	response.Success(c, gin.H{
		"users": data,
		"total": total,
	})
}

// UpdateByUserPasswordToNew Update user's password
// @Summary update user's password to new
// @Description update user's password to new
// @Tags user
// @accept json
// @Produce json
// @Security BearerTokenAuth
// @Param data body types.UpdateUserPasswordRequest true "user information and new password"
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/update [put]
func (h *userHandler) UpdateByUserPasswordToNew(c *gin.Context) {
	form := &types.UpdateUserPasswordRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	// 查询用户是否存在
	user, err := h.iDao.QueryUserByUsername(c, form.Username)
	if err != nil {
		logger.Error("User Not Found", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.AccessDenied)
		return
	}
	logger.Info(form.OldPassword)
	logger.Info(user.Password)
	// 对比用户输入的旧密码与数据库内的旧密码
	err = crypt.CheckSaltPwd(form.OldPassword, user.Password)
	if err != nil {
		logger.Error("AccessDenied", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.AccessDenied)
		return
	}
	// 给新密码加盐
	cryptedNewPwd, err := crypt.GenerateSaltPwd(form.NewPassword)
	if err != nil {
		logger.Error("AccessDenied", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized)
		return
	}
	// 将加盐后的新密码存入数据库
	err = h.iDao.UpdateByIDPasswordToNew(c, user.ID, cryptedNewPwd)
	if err != nil {
		logger.Error("ErrUpdateUser", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrUpdateUser)
		return
	}
	data := &types.UpdateUserPasswordResponse{
		Username: form.Username,
		Password: form.NewPassword,
	}
	response.Success(c, gin.H{"user": data})
}

// LoginUser Update user's password
// @Summary login user
// @Description login user
// @Tags user
// @accept json
// @Produce json
// @Param data body types.LoginUserRequest true "user information and new password"
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/login [post]
func (h *userHandler) LoginUser(c *gin.Context) {
	form := &types.LoginUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	user, err := h.iDao.QueryUserByUsername(c, form.Username)
	if err != nil {
		logger.Error("AccessDenied", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.ErrGetUser)
		return
	}
	err = crypt.CheckSaltPwd(form.Password, user.Password)
	logger.Info(utils.Uint64ToStr(user.ID))
	logger.Info(form.Username)
	logger.Info(form.Password)
	logger.Info(user.Password)
	if err != nil {
		logger.Error("Unauthorized", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.AccessDenied)
		return
	}
	token, err := jwt.GenerateToken(utils.Uint64ToStr(user.ID), form.Username)
	if err != nil {
		logger.Error("Token Error", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.AccessDenied)
		return
	}
	response.Success(c, gin.H{"uid": user.ID, "token": token})
}

// GetUserInfo get a record by uid
// @Summary get user details
// @Description get user details by uid
// @Tags user
// @Security BearerTokenAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/info [get]
func (h *userHandler) GetUserInfo(c *gin.Context) {
	uidAny, exist := c.Get("uid")
	if !exist {
		logger.Error("GetUID error", middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized)
	}
	uidStr := uidAny.(string)
	uid, err := utils.StrToUint64E(uidStr)
	if err != nil {
		logger.Error("GetUID error", logger.Err(err), logger.Any("uid", uid), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.Unauthorized)
	}

	user, err := h.iDao.GetByID(c.Request.Context(), uid)
	if err != nil {
		if errors.Is(err, query.ErrNotFound) {
			logger.Warn("GetUserInfo not found", logger.Err(err), logger.Any("uid", uid), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.Unauthorized)
		} else {
			logger.Error("GetUserInfo error", logger.Err(err), logger.Any("uid", uid), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.Unauthorized)
		}
		return
	}
	data := &types.GetUserByIDRespond{}
	err = copier.Copy(data, user)
	if err != nil {
		response.Error(c, ecode.Unauthorized)
		return
	}
	data.ID = uidStr
	logger.Info(data.ID)
	response.Success(c, gin.H{"user": data})
}

// LogoutUser delete c.Context's key uid
// @Summary delete c.Context's key uid
// @Description delete c.Context's key uid
// @Tags user
// @Security BearerTokenAuth
// @Accept json
// @Produce json
// @Success 200 {object} types.Result{}
// @Router /api/v1/user/logout [post]
func (h *userHandler) LogoutUser(c *gin.Context) {
	keys := c.Keys
	uid, exist := c.Get("uid")
	if exist {
		delete(keys, "uid")
		logger.Info("in delete")
		response.Success(c, gin.H{"uid": uid.(string)})
	}
}

func getUserIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return "", 0, true
	}

	return idStr, id, false
}

func convertUsers(fromValues []*model.User) ([]*types.GetUserByIDRespond, error) {
	toValues := []*types.GetUserByIDRespond{}
	for _, v := range fromValues {
		data := &types.GetUserByIDRespond{}
		err := copier.Copy(data, v)
		if err != nil {
			return nil, err
		}
		data.ID = utils.Uint64ToStr(v.ID)
		toValues = append(toValues, data)
	}

	return toValues, nil
}
