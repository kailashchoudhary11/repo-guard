package initializers

import "os"

func LoadDatabaseURl() string {
	databaseUrl, isPresent := os.LookupEnv("GITHUB_ACCESS_TOKEN")

	if !isPresent {
		return "postgres://localhost:5432/my_db"
	}

	return databaseUrl
}
