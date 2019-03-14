/*
 * Pipeline API
 *
 * Pipeline v0.3.0 swagger
 *
 * API version: 0.3.0
 * Contact: info@banzaicloud.com
 */

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package client

type RestoreResultWarnings struct {
	Ark        []string                            `json:"ark,omitempty"`
	Cluster    []string                            `json:"cluster,omitempty"`
	Namespaces []map[string]map[string]interface{} `json:"namespaces,omitempty"`
}
