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

type LaunchSpotguidesRequest struct {
	RepoOrganization string                 `json:"repoOrganization,omitempty"`
	RepoName         string                 `json:"repoName,omitempty"`
	SpotguideName    string                 `json:"spotguideName,omitempty"`
	Secrets          []CreateSecretRequest  `json:"secrets,omitempty"`
	Values           map[string]interface{} `json:"values,omitempty"`
}
