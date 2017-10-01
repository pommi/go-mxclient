package mxclient

import (
	_ "errors"
	_ "fmt"
)

type MxObject struct {
	Attributes []MxObjectAttribute `json:"attributes"`
	Guid       string              `json:"guid"`
	ObjectType string              `json:"objectType"`
}

type MxObjectAttribute struct {
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
	ReadOnly bool        `json:"readonly,omitempty"`
}

func parseMxObjects(xas_response map[string]interface{}) []MxObject {
	var mxobjects []MxObject

	val, ok := xas_response["mxobjects"]
	if !ok && val == nil {
		return mxobjects
	}

	for _, mxobject_data := range xas_response["mxobjects"].([]interface{}) {
		mxobj, _ := mxobject_data.(map[string]interface{})

		var attributes []MxObjectAttribute
		for attribute, details := range mxobj["attributes"].(map[string]interface{}) {
			attributes = append(attributes, MxObjectAttribute{
				Name:  attribute,
				Value: details.(map[string]interface{})["value"],
				//ReadOnly: details.(map[string]interface{})["readonly"].(bool),
			})
		}

		mxobjects = append(mxobjects, MxObject{
			Attributes: attributes,
			Guid:       mxobj["guid"].(string),
			ObjectType: mxobj["objectType"].(string),
		})
	}

	return mxobjects
}
