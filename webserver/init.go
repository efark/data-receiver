/*
Package webserver has the logic to process the requests.
*/
package webserver

import (
	"fmt"
	"github.com/efark/data-receiver/authenticator"
	"github.com/efark/data-receiver/configuration"
	"github.com/efark/data-receiver/extractor"
	"github.com/efark/data-receiver/logger"
	"github.com/efark/data-receiver/writer"
)

var (
	log, slog = logger.GetLogger()
	services  = make(map[string]*service)
)

type service struct {
	ext  extractor.Extractor
	auth authenticator.Authenticator
	w    writer.Writer
}

// Initialize reads and parses the configuration and stores it in memory.
func Initialize(filepath, cfgInline string) error {

	conf := configuration.NewServiceMap()

	parser := configuration.CreateParser(filepath, cfgInline)
	err := parser.Parse(conf)
	if err != nil {
		slog.Error(err)
		return err
	}

	// log.Info(fmt.Sprintf("%s", conf.List()))
	for _, s := range conf.List() {
		serv, err := conf.Get(s)
		if err != nil {
			slog.Error(err)
			continue
		}
		newExt, err := extractor.CreateExtractor(serv.ExtConfig.Class, serv.ExtConfig.Parameters)
		if err != nil {
			slog.Error(err)
			log.Info(fmt.Sprintf("Extractor for service %q couldn't be created.", s))
			continue
		}

		newAuth, err := authenticator.CreateAuthenticator(serv.AuthConfig.Class, serv.AuthConfig.Parameters)
		if err != nil {
			slog.Error(err)
			log.Info(fmt.Sprintf("Authenticator for service %q couldn't be created.", s))
			continue
		}

		newWriter, err := writer.CreateWriter(serv.WriConfig.Class, serv.WriConfig.Parameters)
		if err != nil {
			slog.Error(err)
			log.Info(fmt.Sprintf("Writer for service %q couldn't be created.", s))
			continue
		}
		SetService(s, newExt, newAuth, newWriter)
	}

	//log.Info(fmt.Sprintf("%+v\n", services))
	return nil
}

// CloseWriters close all the writers for a graceful shutdown.
func CloseWriters() {
	for k, v := range services {
		log.Info(fmt.Sprintf("Closing writer for client: %s.", k))
		v.w.Close()
	}
	return
}

// SetService creates a service with the received name, extractor, authenticator and writer.
func SetService(key string, ext extractor.Extractor, auth authenticator.Authenticator, writer writer.Writer) {
	services[key] = &service{ext: ext, auth: auth, w: writer}
}
