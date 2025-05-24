package utils

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() *zap.Logger{
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
	return Logger
}
