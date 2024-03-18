package cmd

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"hiveify-core/internal/controller/user"
	"hiveify-core/internal/middleware"
)

var (
	// Main 主启动命令
	Main = gcmd.Command{
		Name:  "hiveify-core",
		Usage: "hiveify-core",
		Brief: "hiveify网盘后端服务",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			s := g.Server()
			s.GetOpenApi().Info.Title = "Hiveify-core"
			s.GetOpenApi().Info.Description = "Hiveify网盘后端服务"
			s.GetOpenApi().Config.CommonResponse = ghttp.DefaultHandlerResponse{}
			s.GetOpenApi().Config.CommonResponseDataField = "Data"
			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(ghttp.MiddlewareHandlerResponse, middleware.Auth)
				group.Bind(user.NewV1())
			})
			s.Run()
			return nil
		},
	}
)
