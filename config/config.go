package config

import (
	"flag"
	"os"
)

type Config struct {
	apiPort   string
	secretJWT string
}

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.apiPort, "apiPort", os.Getenv("API_PORT"), "API Port")
	flag.StringVar(&conf.secretJWT, "secretjwt", os.Getenv("SECRET_JWT"), "JWT secret key")
	flag.Parse()

	return conf
}

// func (c *Config) GetDBConnStr() string {
// 	return c.getDBConnStr(c.dbHost, c.dbName)
// }

// func (c *Config) getDBConnStr(dbhost, dbname string) string {
// 	return fmt.Sprintf(
// 		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
// 		c.dbUser,
// 		c.dbPswd,
// 		dbhost,
// 		c.dbPort,
// 		dbname,
// 	)
// }

func (c *Config) GetAPIPort() string {
	return ":" + c.apiPort
}

func (c *Config) GetJWTSecret() []byte {
	return []byte(c.secretJWT)
}
