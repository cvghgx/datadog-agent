// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

//go:build !linux

package agent

import (
	"github.com/DataDog/datadog-agent/comp/core/config"
	log "github.com/DataDog/datadog-agent/comp/core/log/def"
	"github.com/DataDog/datadog-agent/comp/process/types"
	pkgconfigmodel "github.com/DataDog/datadog-agent/pkg/config/model"
	"github.com/DataDog/datadog-agent/pkg/util/flavor"
)

// Enabled determines whether the process agent is enabled based on the configuration
// The process-agent always runs as a stand-alone agent in all non-linux platforms
func Enabled(_ config.Component, _ []types.CheckComponent, _ log.Component) bool {
	return flavor.GetFlavor() == flavor.ProcessAgent
}

// OverrideRunInCoreAgentConfig sets the process_config.run_in_core_agent.enabled to false in unsupported environments.
func OverrideRunInCoreAgentConfig(config config.Component) {
	config.Set("process_config.run_in_core_agent.enabled", false, pkgconfigmodel.SourceAgentRuntime)
}
