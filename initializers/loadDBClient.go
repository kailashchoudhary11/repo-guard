package initializers

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetDBClient(ctx context.Context) (*pgx.Conn, error) {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	dsn := "postgres://" + username + ":" + password + "@" + host + "/" + db_name
	return pgx.Connect(ctx, dsn)
}
