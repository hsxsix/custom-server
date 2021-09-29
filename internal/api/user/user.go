/**
 * @File: user.go
 * @Author: hsien
 * @Description:
 * @Date: 9/18/21 5:03 PM
 */

package user

import (
	"custom_server/internal/consts"
	"custom_server/internal/model/request"
	"custom_server/internal/model/response"
	"custom_server/internal/service/user"
	"custom_server/pkg/log"
	"custom_server/pkg/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type User struct {
	srv *user.Service
}

func NewUser(db *gorm.DB) *User {
	return &User{srv: user.NewService(db)}
}

func (u *User) UserList(ctx *gin.Context) interface{} {
	var query = request.DefaultUserQuery()
	if err := ctx.BindQuery(query); err != nil {
		log.ErrorWithGinCtx(ctx, "bind query params", log.NameError("error", err))
		return consts.BindErrWithError(err)
	}

	res, err := u.srv.GetUserList(util.CtxWithRequestID(ctx), query)
	if err != nil {
		log.ErrorWithGinCtx(ctx, "query user error", log.Any("params", query),
			log.NameError("error", err))
		return err
	}

	return res
}

func (u *User) UserInfo(ctx *gin.Context) interface{} {
	userId := ctx.Query("user_id")
	if userId == "" {
		return nil
	}

	result, err := u.srv.GetUserInfoByID(util.CtxWithRequestID(ctx), userId)
	if err != nil {
		log.ErrorWithGinCtx(ctx, "get user info error",
			log.String("user id", userId), log.NameError("error", err))
		return err
	}
	return result
}

func (u *User) Login(ctx *gin.Context) interface{} {
	login := &request.UserLogin{}
	if err := ctx.ShouldBind(login); err != nil {
		log.ErrorWithGinCtx(ctx, "bind query params", log.NameError("error", err))
		return consts.BindErrWithError(err)
	}

	token, err := u.srv.Login(util.CtxWithRequestID(ctx), login)
	if err != nil {
		log.ErrorWithGinCtx(ctx, "user login failed", log.NameError("error", err))
		return err
	}
	return response.UserTokenRes{
		Token: token,
	}
}

func (u *User) Register(ctx *gin.Context) interface{} {
	register := &request.UserRegister{}
	if err := ctx.ShouldBind(register); err != nil {
		log.ErrorWithGinCtx(ctx, "bind query params", log.NameError("error", err))
		return consts.BindErrWithError(err)
	}

	if err := u.srv.Register(util.CtxWithRequestID(ctx), register); err != nil {
		log.ErrorWithGinCtx(ctx, "register registration failed", log.NameError("error", err))
		return err
	}

	return nil
}
