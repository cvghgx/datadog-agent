// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build !docker

// Package ecs provides the ecs colletor for workloadmeta
package ecs

import (
	wmcatalog "github.com/DataDog/datadog-agent/comp/core/wmcatalog/def"
)

type dependencies struct{}

// NewCollector is a no-op constructor
func NewCollector(deps dependencies) (wmcatalog.Collector, error) {
	return nil, nil
}
