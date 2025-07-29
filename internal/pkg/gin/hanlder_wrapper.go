package gin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/errs"
	"github.com/gin-gonic/gin"
)

// W 封装最基础的 gin.handlerFunc。
func W(bizFunc func(*gin.Context) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		r, err := bizFunc(ctx)
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("unauthorized", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// B 封装从请求体获取参数的 gin.HandlerFunc。
func B[Req any](bizFunc func(*gin.Context, Req) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindJSON(&req); err != nil {
			slog.Error("failed to bind request", slog.Any("err", err))
			return
		}

		r, err := bizFunc(ctx, req)
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("unauthorized", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// Q 封装从 url query 上获取参数的 gin.HandlerFunc。
func Q[Req any](bizFunc func(*gin.Context, Req) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindQuery(&req); err != nil {
			slog.Error("failed to bind request from query", slog.Any("err", err))
			return
		}

		r, err := bizFunc(ctx, req)
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("unauthorized", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// P 封装从 url path 上获取参数的 gin.HandlerFunc。
func P[Req any](bizFunc func(*gin.Context, Req) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindUri(&req); err != nil {
			slog.Error("failed to bind request from uri", slog.Any("err", err))
			return
		}

		r, err := bizFunc(ctx, req)
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("unauthorized", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// WU 封装包含用户登录信息的 gin.handlerFunc。
func WU(bizFunc func(*gin.Context, AuthUser) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		rawVal, ok := ctx.Get(ContextKeyAuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 注意 gin.Context 内的值不能是 *AuthUser
		au, ok := rawVal.(AuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		
		r, err := bizFunc(ctx, au)
		if errors.Is(err, errs.ErrUnauthorized) {
			slog.Debug("unauthorized", slog.Any("err", err))
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// BU 封装从请求体获取参数的且包含用户登录信息 gin.HandlerFunc。
func BU[Req any](bizFunc func(*gin.Context, Req, AuthUser) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindJSON(&req); err != nil {
			slog.Error("failed to bind request", slog.Any("err", err))
			return
		}
		rawVal, ok := ctx.Get(ContextKeyAuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 注意 gin.Context 内的值不能是 *AuthUser
		au, ok := rawVal.(AuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		r, err := bizFunc(ctx, req, au)
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// QU 封装从 url query 上获取参数且包含用户登录信息 gin.HandlerFunc。
func QU[Req any](bizFunc func(*gin.Context, Req, AuthUser) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindQuery(&req); err != nil {
			slog.Error("failed to bind request from query", slog.Any("err", err))
			return
		}
		rawVal, ok := ctx.Get(ContextKeyAuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 注意 gin.Context 内的值不能是 *AuthUser
		au, ok := rawVal.(AuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		r, err := bizFunc(ctx, req, au)
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}

// PU 封装从 url path 上获取参数且包含用户登录信息 gin.HandlerFunc。
func PU[Req any](bizFunc func(*gin.Context, Req, AuthUser) (R, error)) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		if err := ctx.BindUri(&req); err != nil {
			slog.Error("failed to bind request from uri", slog.Any("err", err))
			return
		}
		rawVal, ok := ctx.Get(ContextKeyAuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 注意 gin.Context 内的值不能是 *AuthUser
		au, ok := rawVal.(AuthUser)
		if !ok {
			slog.Error("failed to get auth user")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		r, err := bizFunc(ctx, req, au)
		if err != nil {
			slog.Error("failed to handle request", slog.Any("err", err))
			ctx.PureJSON(http.StatusInternalServerError, r)
			return
		}
		ctx.PureJSON(r.Code, r)
	}
}
