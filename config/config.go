package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/onunez-g/auth-api/mail"
)

type Config struct {
	dbUser         string
	dbPswd         string
	dbHost         string
	dbPort         string
	dbName         string
	apiPort        string
	secretJWT      string
	githubClientId string
	githubSecretId string
	sessionSecret  string
	emailUser      string
	emailPass      string
	emailHost      string
	emailPort      string
}

var Cfg *Config

func Get() *Config {
	conf := &Config{}

	flag.StringVar(&conf.dbUser, "dbuser", os.Getenv("POSTGRES_USER"), "DB user name")
	flag.StringVar(&conf.dbPswd, "dbpswd", os.Getenv("POSTGRES_PASSWORD"), "DB pass")
	flag.StringVar(&conf.dbPort, "dbport", os.Getenv("POSTGRES_PORT"), "DB port")
	flag.StringVar(&conf.dbHost, "dbhost", os.Getenv("POSTGRES_HOST"), "DB host")
	flag.StringVar(&conf.dbName, "dbname", os.Getenv("POSTGRES_DB"), "DB name")
	flag.StringVar(&conf.apiPort, "apiPort", os.Getenv("API_PORT"), "API Port")
	flag.StringVar(&conf.secretJWT, "secretjwt", os.Getenv("SECRET_JWT"), "JWT secret key")
	flag.StringVar(&conf.githubClientId, "githubClientId", os.Getenv("CLIENT_ID"), "Github client ID")
	flag.StringVar(&conf.githubSecretId, "githubSecretId", os.Getenv("SECRET_ID"), "Github secret ID")
	flag.StringVar(&conf.sessionSecret, "sessionSecret", os.Getenv("SESSION"), "Cookie session secret")
	flag.StringVar(&conf.emailUser, "emailUser", os.Getenv("SMTP_USER"), "smtp server user")
	flag.StringVar(&conf.emailPass, "emailPass", os.Getenv("SMTP_PASS"), "smtp server password")
	flag.StringVar(&conf.emailHost, "emailHost", os.Getenv("SMTP_HOST"), "smtp server host")
	flag.StringVar(&conf.emailPort, "emailPort", os.Getenv("SMTP_PORT"), "smtp server port")
	flag.Parse()

	return conf
}

func (c *Config) GetDBConnStr() string {
	return c.getDBConnStr(c.dbHost, c.dbName)
}

func (c *Config) getDBConnStr(dbhost, dbname string) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.dbUser,
		c.dbPswd,
		dbhost,
		c.dbPort,
		dbname,
	)
}

func (c *Config) GetAPIPort() string {
	return ":" + c.apiPort
}

func (c *Config) GetJWTSecret() []byte {
	return []byte(c.secretJWT)
}

func (c *Config) GetGithubClientId() string {
	return c.githubClientId
}

func (c *Config) GetGithubSecretId() string {
	return c.githubSecretId
}

func (c *Config) GetSessionSecret() []byte {
	return []byte(c.sessionSecret)
}

func (c *Config) GetSMTPUser() string {
	return c.emailUser
}

func (c *Config) GetSMTPPassword() string {
	return c.emailPass
}

func (c *Config) GetSMTPSettings(to string) mail.MailSettings {
	return mail.MailSettings{
		User: c.emailUser,
		Pass: c.emailPass,
		From: c.emailUser,
		Smtp: c.emailHost,
		Port: c.emailPort,
		To:   to,
	}
}
