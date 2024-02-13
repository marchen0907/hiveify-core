package main

import (
	"github.com/charmbracelet/log"
	"github.com/joho/godotenv"
	"hiveify-core/database/postgresql"
	"hiveify-core/routers"
	"os"
)

func loadEnv(envFile string) {
	err := godotenv.Load(envFile)
	if err != nil {
		panic("Error loading env file")
	}
}

func main() {
	// 获取编译时传递的环境变量文件路径
	envFile := os.Getenv("ENV_FILE")

	// 加载环境变量文件
	loadEnv(envFile)

	defer postgresql.DBConnection().Close()

	router := routers.Router()
	err := router.Run("0.0.0.0:8080")
	if err != nil {
		log.Errorf("Failed to run server: %s", err.Error())
		panic(err)
	}
}
