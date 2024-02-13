package postgresql

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"
	"os"
)

// Ctx 数据库的初始化上下文
var (
	Ctx            = context.Background()
	postgreSQLPool *pgxpool.Pool
)

// DBConnection 连接数据库
func DBConnection() *pgxpool.Pool {
	if postgreSQLPool == nil {
		connStr := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
			os.Getenv("POSTGRESQL_USER"),
			url.QueryEscape(os.Getenv("POSTGRESQL_PASSWORD")),
			os.Getenv("POSTGRESQL_HOST"), os.Getenv("POSTGRESQL_PORT"), os.Getenv("POSTGRESQL_DATABASE"))
		var err error
		postgreSQLPool, err = pgxpool.New(context.Background(), connStr)
		if err != nil {
			log.Errorf("Unable to create connection pool: %s\n", err.Error())
			os.Exit(1)
		}
	}
	return postgreSQLPool
}
