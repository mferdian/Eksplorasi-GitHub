package middleware

import (
	"net/http"
	"strings"

	"github.com/Amierza/go-boiler-plate/dto"
	"github.com/Amierza/go-boiler-plate/jwt"
	"github.com/Amierza/go-boiler-plate/response"
	"github.com/gin-gonic/gin"
)

func Authentication(jwtService jwt.IJWTService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_FOUND, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !strings.Contains(authHeader, "Bearer") {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		authHeader = strings.Replace(authHeader, "Bearer ", "", -1)
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_NOT_VALID, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		if !token.Valid {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, dto.MESSAGE_FAILED_TOKEN_DENIED_ACCESS, nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		userID, err := jwtService.GetUserIDByToken(authHeader)
		if err != nil {
			res := response.BuildResponseFailed(dto.MESSAGE_FAILED_PROSES_REQUEST, err.Error(), nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, res)
			return
		}

		ctx.Set("Authorization", authHeader)
		ctx.Set("user_id", userID)
		ctx.Next()
	}
}
