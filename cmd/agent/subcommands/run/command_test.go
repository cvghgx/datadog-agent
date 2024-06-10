// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package run

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/DataDog/datadog-agent/cmd/agent/command"
	"github.com/DataDog/datadog-agent/comp/core"
	"github.com/DataDog/datadog-agent/comp/core/pid/pidimpl"
	"github.com/DataDog/datadog-agent/comp/core/secrets"
	"github.com/DataDog/datadog-agent/pkg/util/fxutil"
)

func TestCommand(t *testing.T) {
	fxutil.TestOneShotSubcommand(t,
		Commands(newGlobalParamsTest(t)),
		[]string{"run"},
		run,
		func(pidParam pidimpl.Params, coreParams core.BundleParams, secretParams secrets.Params) {
			require.Equal(t, true, secretParams.Enabled)
		})
}

func TestCommandPidfile(t *testing.T) {
	fxutil.TestOneShotSubcommand(t,
		Commands(newGlobalParamsTest(t)),
		[]string{"run", "--pidfile", "/pid/file"},
		run,
		func(pidParams pidimpl.Params, coreParams core.BundleParams, secretParams secrets.Params) {
			require.Equal(t, "/pid/file", pidParams.PIDfilePath)
			require.Equal(t, true, secretParams.Enabled)
		})
}

func newGlobalParamsTest(t *testing.T) *command.GlobalParams {
	// Because getSharedFxOption uses fx.Invoke, demultiplexer component is built
	// which lead to build:
	//   - configPath.Component which requires a valid datadog.yaml
	//   - hostname.Component which requires a valid hostname
	tempDirPath := t.TempDir()
	configPath := path.Join(tempDirPath, "datadog.yaml")
	configContent := fmt.Sprintf(`
hostname: test
run_path: %s
`, tempDirPath)

	err := os.WriteFile(configPath, []byte(configContent), 0o644)
	require.NoError(t, err)

	return &command.GlobalParams{
		ConfFilePath: configPath,
	}
}
