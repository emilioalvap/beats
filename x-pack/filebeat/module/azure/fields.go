// Copyright Elasticsearch B.V. and/or licensed to Elasticsearch B.V. under one
// or more contributor license agreements. Licensed under the Elastic License;
// you may not use this file except in compliance with the Elastic License.

// Code generated by beats/dev-tools/cmd/asset/asset.go - DO NOT EDIT.

package azure

import (
	"github.com/elastic/beats/v7/libbeat/asset"
)

func init() {
	if err := asset.SetFields("filebeat", "azure", asset.ModuleFieldsPri, AssetAzure); err != nil {
		panic(err)
	}
}

// AssetAzure returns asset data.
// This is the base64 encoded zlib format compressed contents of module/azure.
func AssetAzure() string {
	return "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"
}
