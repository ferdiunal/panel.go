package constants

import "os"

var (
	APP_URL = os.Getenv("APP_URL")
	APP_ENV = os.Getenv("APP_ENV")
)
