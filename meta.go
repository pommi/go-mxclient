package mxclient

import (
	"errors"
	"fmt"
)

type Metadata struct {
	Entities []Entity
}

type Entity struct {
	Name        string
	Persistable bool
	Attributes  []Attribute
}

type Attribute struct {
	Name      string
	Type      string
	Reference string
}

func (c *Client) GetMetadata() (metadata *Metadata, err error) {
	var entities []Entity

	session_data, err := c.GetSessionData()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while retrieving session data: %s", err))
	}

	for _, entity_data := range session_data["metadata"].([]interface{}) {
		e, _ := entity_data.(map[string]interface{})

		var attributes []Attribute
		for attribute, details := range e["attributes"].(map[string]interface{}) {
			if details.(map[string]interface{})["type"] == "ObjectReference" {
				attributes = append(attributes, Attribute{
					Name:      attribute,
					Type:      details.(map[string]interface{})["type"].(string),
					Reference: details.(map[string]interface{})["klass"].(string),
				})
			} else {
				attributes = append(attributes, Attribute{
					Name: attribute,
					Type: details.(map[string]interface{})["type"].(string),
				})
			}
		}

		entities = append(entities, Entity{
			Name:        e["objectType"].(string),
			Persistable: e["persistable"].(bool),
			Attributes:  attributes,
		})
	}
	return &Metadata{
		Entities: entities,
	}, nil
}
