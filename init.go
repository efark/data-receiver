/*
This file initializes the logger for the 'main' package.
*/
package main

import (
	"github.com/efark/data-receiver/logger"
	"go.uber.org/zap"
)

var (
	log  *zap.Logger
	slog *zap.SugaredLogger
)

func init() {
	log, slog = logger.GetLogger()
}
