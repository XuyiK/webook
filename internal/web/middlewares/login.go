package middlewares

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type LoginMiddlewareBuilder struct {
}

func (m *LoginMiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	// 注册一下这个类型
	gob.Register(time.Now())
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/users/signup" || path == "/users/login" {
			// 不需要登录校验
			return
		}
		sess := sessions.Default(ctx)
		const userIdKey = "userId"
		userId := sess.Get(userIdKey)
		if userId == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		now := time.Now()
		const updateTimeKey = "update_time"
		val := sess.Get(updateTimeKey)
		lastUpdateTime, ok := val.(time.Time)
		if val == nil || (!ok) || now.Sub(lastUpdateTime) > time.Minute {
			sess.Set(updateTimeKey, now)
			// 由于sess.Set()是覆盖式更新，所以需要重新set userId
			sess.Set(userIdKey, userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}

	}
}
