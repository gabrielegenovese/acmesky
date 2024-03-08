package ginContextRepo

import (
	"sync"

	"github.com/gin-gonic/gin"
)

var ginContextStore sync.Map

func SetContext(id string, ctx *gin.Context) {
	var key string = id
	ginContextStore.Store(key, ctx)
}

func GetContext(id string) *gin.Context {

	var key string = id

	ctx, _ := ginContextStore.Load(key)
	if ctx != nil {
		return ctx.(*gin.Context)
	}
	return nil
}

func UnsetContext(id string) {
	var key string = id
	ginContextStore.Delete(key)
}
