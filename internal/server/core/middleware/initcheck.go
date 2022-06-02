package middleware

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func InitCheck() iris.Handler {
	return func(ctx *context.Context) {
		//lang := ctx.URLParam("lang")
		//if lang != commConsts.Language {
		//	commConsts.Language = lang
		//	i118Utils.Init(commConsts.Language, commConsts.AppServer)
		//}

		ctx.Next()
	}
}
