package gin

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/JrMarcco/kuryr-admin/internal/errs"
	"github.com/gin-gonic/gin"
)

func W(bizFunc func(ctx *gin.Context) (R, error)) gin.HandlerFunc {
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
		ctx.PureJSON(http.StatusOK, r)
	}
}

func B[Req any](bizFunc func(ctx *gin.Context, req Req) (R, error)) gin.HandlerFunc {
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
		ctx.PureJSON(http.StatusOK, r)
	}
}
