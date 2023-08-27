package api

import (
	"errors"
	"strings"

	"github.com/ShenPingYuan/go-webdemo/token"
	"github.com/gin-gonic/gin"
)

const (
	authHeaderKey  = "Authorization"
	authTypeBearer = "bearer"
	currentUserKey = "user"
)

func authMiddleware(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationToken := ctx.GetHeader(authHeaderKey)
		if len(authorizationToken) == 0 {
			err := errors.New("authorization header is not provided")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		fields := strings.Fields(authorizationToken)
		if len(fields) < 2 {
			err := errors.New("invalid authorization header format")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authTypeBearer {
			err := errors.New("unsupported authorization type")
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}

		payload, err := tokenMaker.VerifyToken(fields[1])
		if err != nil {
			ctx.AbortWithStatusJSON(401, errorResponse(err))
			return
		}
		ctx.Set(currentUserKey, payload)
		ctx.Next()
	}
}
