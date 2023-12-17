package springcloud

import (
	"errors"
)

var (
	EmptyPropertySource = errors.New("no configurations in property source")
	SourceHasNoConfigs  = errors.New("source has no configurations inside")
)

// SpringResponse Structs having same structure as want from Cloud config
type SpringResponse struct {
	Name            string           `json:"name"`
	Profiles        []string         `json:"profiles"`
	Label           string           `json:"label"`
	Version         string           `json:"version"`
	PropertySources []propertySource `json:"propertySources"`
}

type propertySource struct {
	Name   string                 `json:"name"`
	Source map[string]interface{} `json:"source"`
}

func (cc *SpringResponse) ToMap() (map[string]interface{}, error) {
	if len(cc.PropertySources) < 1 {
		return nil, EmptyPropertySource
	}
	if cc.PropertySources[0].Source == nil {
		return nil, SourceHasNoConfigs
	}
	return cc.PropertySources[0].Source, nil
}
