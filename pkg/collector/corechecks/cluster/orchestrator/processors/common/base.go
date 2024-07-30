// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//go:build orchestrator

// Package common provides basic handlers used by orchestrator processor
package common

import (
	model "github.com/DataDog/agent-payload/v5/process"
	"github.com/DataDog/datadog-agent/pkg/collector/corechecks/cluster/orchestrator/processors"
	"github.com/DataDog/datadog-agent/pkg/util/installinfo"
	"github.com/DataDog/datadog-agent/pkg/version"
)

//nolint:revive // TODO(CAPP) Fix revive linter
type BaseHandlers struct{}

//nolint:revive // TODO(CAPP) Fix revive linter
func (BaseHandlers) BeforeCacheCheck(ctx processors.ProcessorContext, resource, resourceModel interface{}) (skip bool) {
	return
}

//nolint:revive // TODO(CAPP) Fix revive linter
func (BaseHandlers) BeforeMarshalling(ctx processors.ProcessorContext, resource, resourceModel interface{}) (skip bool) {
	return
}

//nolint:revive // TODO(CAPP) Fix revive linter
func (BaseHandlers) ScrubBeforeMarshalling(ctx processors.ProcessorContext, resource interface{}) {}

//nolint:revive // TODO(CAPP) Fix revive linter
func (BaseHandlers) BuildManifestMessageBody(ctx processors.ProcessorContext, resourceManifests []interface{}, groupSize int) model.MessageBody {
	return ExtractModelManifests(ctx, resourceManifests, groupSize)
}

// ExtractModelManifests creates the model manifest from the given manifests
func ExtractModelManifests(ctx processors.ProcessorContext, resourceManifests []interface{}, groupSize int) *model.CollectorManifest {
	pctx := ctx.(*processors.K8sProcessorContext)
	manifests := make([]*model.Manifest, 0, len(resourceManifests))

	for _, m := range resourceManifests {
		manifests = append(manifests, m.(*model.Manifest))
	}

	cm := &model.CollectorManifest{
		ClusterName: pctx.Cfg.KubeClusterName,
		ClusterId:   pctx.ClusterID,
		Manifests:   manifests,
		GroupId:     pctx.MsgGroupID,
		GroupSize:   int32(groupSize),
		Tags:        append(pctx.Cfg.ExtraTags, pctx.ApiGroupVersionTag),
	}
	return cm
}

// K8sAgentInfo returns the agent info for k8s
func K8sAgentInfo(i *installinfo.InstallInfo) *model.K8SAgentInfo {
	if i == nil {
		return &model.K8SAgentInfo{
			Version: version.AgentVersion,
		}
	}
	return &model.K8SAgentInfo{
		Version:              version.AgentVersion,
		InstallMethod:        i.Tool,
		InstallMethodVersion: i.InstallerVersion,
	}
}
