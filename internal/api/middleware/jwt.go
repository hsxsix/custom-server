/**
 * @File: jwt.go
 * @Author: hsien
 * @Description:
 * @Date: 9/24/21 11:41 PM
 */

package middleware

import (
	"custom_server/internal/consts"
	"custom_server/internal/model"
	"custom_server/internal/model/response"
	"custom_server/pkg/errno"
	"custom_server/pkg/log"
	"custom_server/pkg/util"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := func(ctx *gin.Context) *errno.Err {
			token := ctx.Request.Header.Get("Authorization")
			if token == "" {
				return consts.PermissionDenied
			}

			claims, err := util.ParseToken(token, consts.JwtSecret, &model.UserClaims{})
			if err != nil {
				return consts.TokenInvalidWithError(err)
			}
			if err := claims.Valid(); err != nil {
				return consts.TokenExpiredWithError(err)
			}
			return nil
		}(ctx)

		if err != nil {
			log.ErrorWithGinCtx(ctx, "authorization failed", log.NameError("error", err))
			ctx.JSON(http.StatusOK, response.Response{
				Code: err.Code,
				Msg:  err.Msg,
			})
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}

func ParseUserClaims(ctx *gin.Context) (*model.UserClaims, error) {
	token := ctx.Request.Header.Get("Authorization")
	claims, err := util.ParseToken(token, consts.JwtSecret, &model.UserClaims{})
	if err != nil {
		return nil, err
	}

	user, ok := claims.(*model.UserClaims)
	if !ok {
		return nil, errors.New("convert to model.UserClaims failed")
	}
	return user, nil
}
