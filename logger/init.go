/*
Package logger sets up logging for the application, based on Uber zap's logger.
*/
package logger

import (
	"go.uber.org/zap"
	"sync"
)

// Package internal variable to implement singleton
var (
	innerLogger *zap.Logger
	innerSugar  *zap.SugaredLogger
	onceLogger  sync.Once
)

// GetLogger returns singleton logger object.
func GetLogger() (*zap.Logger, *zap.SugaredLogger) {
	var err error

	onceLogger.Do(func() {
		innerLogger, err = zap.NewDevelopment()
		if err != nil {
			panic("Unable to create logger. Quitting application.")
		}
		innerSugar = innerLogger.Sugar()
	})
	return innerLogger, innerSugar
}
