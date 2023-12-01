package config

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadConfig(path string) IConfig {
	envMap, err := godotenv.Read(path)
	if err != nil {
		log.Fatalf("load dotenv failed: %v", err)
	}

	appConfig := &app{
		host: envMap["APP_HOST"],
		port: func() int {
			p, err := strconv.Atoi(envMap["APP_PORT"])
			if err != nil {
				log.Fatalf("load port failed: %v", err)
			}
			return p
		}(),
		name:    envMap["APP_NAME"],
		version: envMap["APP_VERSION"],
		readTimeout: func() time.Duration {
			t, err := strconv.Atoi(envMap["APP_READ_TIMEOUT"])
			if err != nil {
				log.Fatalf("load read timeout failed: %v", err)
			}
			return time.Duration(int64(t) * int64(math.Pow10(9)))
		}(),
		writeTimeout: func() time.Duration {
			t, err := strconv.Atoi(envMap["APP_WRITE_TIMEOUT"])
			if err != nil {
				log.Fatalf("load write timeout failed: %v", err)
			}
			return time.Duration(int64(t) * int64(math.Pow10(9)))
		}(),
		bodyLimit: func() int {
			b, err := strconv.Atoi(envMap["APP_BODY_LIMIT"])
			if err != nil {
				log.Fatalf("load body limit failed: %v", err)
			}
			return b
		}(),
		fileLimit: func() int {
			f, err := strconv.Atoi(envMap["APP_FILE_LIMIT"])
			if err != nil {
				log.Fatalf("load file limit failed: %v", err)
			}
			return f
		}(),
	}

	dbConfig := &db{
		host: envMap["DB_HOST"],
		port: func() int {
			p, err := strconv.Atoi(envMap["DB_PORT"])
			if err != nil {
				log.Fatalf("load db port failed: %v", err)
			}
			return p
		}(),
		protocol: envMap["DB_PROTOCOL"],
		username: envMap["DB_USERNAME"],
		password: envMap["DB_PASSWORD"],
		database: envMap["DB_DATABASE"],
		sslMode:  envMap["DB_SSL_MODE"],
		maxConnections: func() int {
			m, err := strconv.Atoi(envMap["DB_MAX_CONNECTIONS"])
			if err != nil {
				log.Fatalf("load db max connections failed: %v", err)
			}
			return m
		}(),
	}

	return &config{
		app: appConfig,
		db:  dbConfig,
	}
}

type IConfig interface {
	App() IAppConfig
	DB() IDBConfig
}

type config struct {
	app *app
	db  *db
}

type IAppConfig interface {
	Url() string
	Name() string
	Version() string
	Host() string
	Port() int
	ReadTimeout() time.Duration
	WriteTimeout() time.Duration
	BodyLimit() int
	FileLimit() int
}

type app struct {
	host         string
	port         int
	name         string
	version      string
	readTimeout  time.Duration
	writeTimeout time.Duration
	bodyLimit    int // bytes
	fileLimit    int // bytes
}

func (c *config) App() IAppConfig {
	return c.app
}

func (a *app) Url() string                 { return fmt.Sprintf("%s:%d", a.host, a.port) }
func (a *app) Name() string                { return a.name }
func (a *app) Version() string             { return a.version }
func (a *app) Host() string                { return a.host }
func (a *app) Port() int                   { return a.port }
func (a *app) ReadTimeout() time.Duration  { return a.readTimeout }
func (a *app) WriteTimeout() time.Duration { return a.writeTimeout }
func (a *app) BodyLimit() int              { return a.bodyLimit }
func (a *app) FileLimit() int              { return a.fileLimit }

type IDBConfig interface {
	Url() string
	MaxConnections() int
}

type db struct {
	host           string
	port           int
	protocol       string
	username       string
	password       string
	database       string
	sslMode        string
	maxConnections int
}

func (c *config) DB() IDBConfig {
	return c.db
}

func (d *db) Url() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.host,
		d.port,
		d.username,
		d.password,
		d.database,
		d.sslMode,
	)
}

func (d *db) MaxConnections() int { return d.maxConnections }
