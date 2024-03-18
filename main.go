package main

import (
	_ "github.com/gogf/gf/contrib/drivers/pgsql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	_ "hiveify-core/internal/logic"
	_ "hiveify-core/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"hiveify-core/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.GetInitCtx())
}
