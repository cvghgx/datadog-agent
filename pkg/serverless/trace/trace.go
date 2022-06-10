// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package trace

import (
	"context"
	"strings"

	tracecmdconfig "github.com/DataDog/datadog-agent/cmd/trace-agent/config"
	ddConfig "github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/trace/agent"
	"github.com/DataDog/datadog-agent/pkg/trace/config"
	"github.com/DataDog/datadog-agent/pkg/trace/pb"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

// ServerlessTraceAgent represents a trace agent in a serverless context
type ServerlessTraceAgent struct {
	ta           *agent.Agent
	spanModifier *spanModifier
	cancel       context.CancelFunc
}

// Load abstracts the file configuration loading
type Load interface {
	Load() (*config.AgentConfig, error)
}

// LoadConfig is implementing Load to retrieve the config
type LoadConfig struct {
	Path string
}

// httpURLMetaKey is the key of the span meta containing the HTTP URL
const httpURLMetaKey = "http.url"

// lambdaRuntimeUrlPrefix is the first part of a URL for a call to the Lambda runtime API
const lambdaRuntimeURLPrefix = "http://127.0.0.1:9001"

// lambdaExtensionURLPrefix is the first part of a URL for a call from the Datadog Lambda Library to the Lambda Extension
const lambdaExtensionURLPrefix = "http://127.0.0.1:8124"

const invocationSpanResource = "dd-tracer-serverless-span"

// Load loads the config from a file path
func (l *LoadConfig) Load() (*config.AgentConfig, error) {
	return tracecmdconfig.LoadConfigFile(l.Path)
}

// Start starts the agent
func (s *ServerlessTraceAgent) Start(enabled bool, loadConfig Load) {
	if enabled {
		// during hostname resolution the first step is to make a GRPC call which is timeboxed to a 2 seconds deadline
		// in the serverless mode, we don't start the GRPC server so this call will fail and cause a 2 seconds delay
		// by setting cmd_port to -1, this will cause the GRPC client to fail instantly
		ddConfig.Datadog.Set("cmd_port", "-1")

		tc, confErr := loadConfig.Load()
		if confErr != nil {
			log.Errorf("Unable to load trace agent config: %s", confErr)
		} else {
			context, cancel := context.WithCancel(context.Background())
			tc.Hostname = ""
			tc.SynchronousFlushing = true
			s.ta = agent.NewAgent(context, tc)
			s.spanModifier = &spanModifier{}
			s.ta.ModifySpan = s.spanModifier.ModifySpan
			s.ta.DiscardSpan = filterSpanFromLambdaLibraryOrRuntime
			s.cancel = cancel
			go func() {
				s.ta.Run()
			}()
		}
	}
}

// Get returns the trace agent instance
func (s *ServerlessTraceAgent) Get() *agent.Agent {
	return s.ta
}

// SetTags sets the tags to the trace agent config and span processor
func (s *ServerlessTraceAgent) SetTags(tagMap map[string]string) {
	s.ta.SetGlobalTagsUnsafe(tagMap)
	s.spanModifier.tags = tagMap
}

// Stop stops the trace agent
func (s *ServerlessTraceAgent) Stop() {
	if s.cancel != nil {
		s.cancel()
	}
}

// filterSpanFromLambdaLibraryOrRuntime returns true if a span was generated by internal HTTP calls within the Datadog
// Lambda Library or the Lambda runtime
func filterSpanFromLambdaLibraryOrRuntime(span *pb.Span) bool {
	if val, ok := span.Meta[httpURLMetaKey]; ok {
		if strings.HasPrefix(val, lambdaExtensionURLPrefix) {
			log.Debugf("Detected span with http url %s, removing it", val)
			return true
		}
		if strings.HasPrefix(val, lambdaRuntimeURLPrefix) {
			log.Debugf("Detected span with http url %s, removing it", val)
			return true
		}
	}
	if span != nil && span.Resource == invocationSpanResource {
		log.Debugf("Detected invocation span from tracer, removing it")
		return true
	}
	return false
}
