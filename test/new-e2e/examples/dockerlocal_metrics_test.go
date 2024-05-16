// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2024-present Datadog, Inc.

package examples

import (
	"github.com/DataDog/test-infra-definitions/components/datadog/agentparams"
	"testing"
	"time"

	"github.com/DataDog/datadog-agent/test/fakeintake/client"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/e2e"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments"
	"github.com/DataDog/datadog-agent/test/new-e2e/pkg/environments/local"

	"github.com/stretchr/testify/assert"
)

type localFakeintakeSuiteMetrics struct {
	e2e.BaseSuite[environments.DockerLocal]
}

func TestVMSuiteEx5Local(t *testing.T) {
	suiteParams := []e2e.SuiteOption{
		e2e.WithProvisioner(
			local.Provisioner(
				local.WithAgentOptions(
					agentparams.WithLatest(),
					// Setting hostname to test name due to fact Agent can't
					// work out it's hostname in a container correctly
					agentparams.WithHostname(t.Name())))),
	}

	if isDevModeEnabled {
		suiteParams = append(suiteParams, e2e.WithDevMode())
	}

	e2e.Run(t, &localFakeintakeSuiteMetrics{}, suiteParams...)
}

func (v *localFakeintakeSuiteMetrics) Test1_FakeIntakeReceivesMetrics() {
	v.EventuallyWithT(func(c *assert.CollectT) {
		metricNames, err := v.Env().FakeIntake.Client().GetMetricNames()
		assert.NoError(c, err)
		assert.Greater(c, len(metricNames), 0)
	}, 5*time.Minute, 10*time.Second)
}

func (v *localFakeintakeSuiteMetrics) Test2_FakeIntakeReceivesSystemLoadMetric() {
	v.EventuallyWithT(func(c *assert.CollectT) {
		metrics, err := v.Env().FakeIntake.Client().FilterMetrics("system.load.1")
		assert.NoError(c, err)
		assert.Greater(c, len(metrics), 0, "no 'system.load.1' metrics yet")
	}, 5*time.Minute, 10*time.Second)
}

func (v *localFakeintakeSuiteMetrics) Test3_FakeIntakeReceivesSystemUptimeHigherThanZero() {
	v.EventuallyWithT(func(c *assert.CollectT) {
		metrics, err := v.Env().FakeIntake.Client().FilterMetrics("system.uptime", client.WithMetricValueHigherThan(0))
		assert.NoError(c, err)
		assert.Greater(c, len(metrics), 0, "no 'system.uptime' with value higher than 0 yet")
		assert.Greater(c, len(metrics), 0, "no 'system.load.1' metrics yet")
	}, 5*time.Minute, 10*time.Second)
}
