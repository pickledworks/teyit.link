package utils

import (
	"github.com/kelseyhightower/envconfig"
	"log"
	"time"
)

// Config for the application
type Config struct {
	// Server address to bind
	ServerAddr string
	// Which dialect to pass to GORM (defaults to "mysql")
	DbDialect string
	// Database connection URI (defaults to "link:root@/teyitlink", see README)
	DbUri string
	// AWS S3 bucket the archive files are stored in (defaults to "teyitlink")
	BucketName string
	// Which Lambda to use for archive worker
	WorkerLambdaName string
	// Which AWS region to use
	AwsRegion string
	// the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m
	GracefulShutdown time.Duration
}

var Conf *Config

func InitConfig() *Config {
	c := new(Config)
	err := envconfig.Process("", c)
	if err != nil {
		log.Fatal("error while InitConfig", err)
	}
	c.setDefaults()
	Conf = c
	return c
}

func (c *Config) setDefaults() {
	if c.ServerAddr == "" {
		c.ServerAddr = "0.0.0.0:8080"
	}
	if c.DbDialect == "" {
		c.DbDialect = "mysql"
	}
	if c.DbUri == "" {
		c.DbUri = "link:root@/teyitlink"
	}
	if c.WorkerLambdaName == "" {
		c.WorkerLambdaName = "teyitlink-archive"
	}
	if c.BucketName == "" {
		c.BucketName = "teyitlink"
	}
	if c.AwsRegion == "" {
		c.AwsRegion = "eu-central-1"
	}
	if c.GracefulShutdown == 0 {
		c.GracefulShutdown = time.Second * 15
	}
}

func GetConfig() *Config {
	return Conf
}
