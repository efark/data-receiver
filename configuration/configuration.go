/*
Package configuration contains Configuration and Parser interfaces.
This file contains the Configuration interface.
*/
package configuration

import (
	"fmt"
)

// Configuration interface is the main interface of the package.
type Configuration interface {
	Add(string, *ServiceConfig) error
	List() []string
	Get(string) (*ServiceConfig, error)
}

// ServiceMap implements configuration interface, it has a map to hold the configuration for all the services of the webserver.
type ServiceMap struct {
	Services map[string]*ServiceConfig `json:"services" yaml:"services"`
}

// NewServiceMap returns a pointer to a ServiceMap struct, holding an empty map.
func NewServiceMap() *ServiceMap {
	return &ServiceMap{Services: make(map[string]*ServiceConfig)}
}

// Add adds the config for a service to the map in the ServiceMap.
func (c *ServiceMap) Add(name string, conf *ServiceConfig) error {
	if _, ok := c.Services[name]; ok {
		return fmt.Errorf("Service %q already exists.", name)
	}
	c.Services[name] = conf
	return nil
}

// List returns the list of names of the services in the inner map.
func (c *ServiceMap) List() []string {
	temp := make([]string, 0, len(c.Services))
	for k := range c.Services {
		temp = append(temp, k)
	}
	return temp
}

// Get returns the config for a specific Service or an error if the service doesn't exist.
func (c *ServiceMap) Get(name string) (*ServiceConfig, error) {
	s, ok := c.Services[name]
	if !ok {
		err := fmt.Errorf("Service %q not found.", name)
		return s, err
	}
	return s, nil
}

// ServiceConfig has the necessary fields to store the configuration of each service.
type ServiceConfig struct {
	ExtConfig  *SimpleConfig `json:"extractor" yaml:"extractor"`
	AuthConfig *SimpleConfig `json:"authenticator" yaml:"authenticator"`
	WriConfig  *SimpleConfig `json:"writer" yaml:"writer"`
}

// NewServiceConfig generates the config for a service based on the Config for each module.
func NewServiceConfig(ext, auth, w *SimpleConfig) *ServiceConfig {
	return &ServiceConfig{ext, auth, w}
}

// SimpleConfig is a basic config that has a Class field to define the type of module (ie, MemoryWriter for Writer or HeaderExtractor for Header),
// and a map to hold the parameters.
type SimpleConfig struct {
	Class      string            `json:"type" yaml:"type"`
	Parameters map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
}

// NewSimpleConfig creates the config using the type and the parameters received.
func NewSimpleConfig(class string, params map[string]string) *SimpleConfig {
	return &SimpleConfig{class, params}
}
