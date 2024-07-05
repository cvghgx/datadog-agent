// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2023-present Datadog, Inc.

// Package pipeline provides helper functions for working with the Gitlab pipeline
package pipeline

// GetPipelineMSIURL returns the URL for the agent MSI built by the pipeline
import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

const (
	DefaultMajorVersion           = "7"
	DefaultArch                   = "x86_64"
	AgentInstallerListProductName = "datadog-agent"
	AgentS3BucketRelease          = "ddagent-windows-stable"
	AgentS3BucketTesting          = "dd-agent-mstesting"
	BetaChannel                   = "beta"
	BetaURL                       = "https://s3.amazonaws.com/dd-agent-mstesting/builds/beta/installers_v2.json"
	StableChannel                 = "stable"
	StableURL                     = "https://ddagent-windows-stable.s3.amazonaws.com/installers_v2.json"
)

// GetArtifact searches a public S3 bucket for a given artifact from a Gitlab pipeline
// majorVersion = [6,7]
// predicate = A function taking the artifact name (from github.com/aws/aws-sdk-go-v2/service/s3/types.Object.Key)
// and that returns true when the artifact matches.
func GetArtifact(pipelineID, bucket, majorVersion string, predicate func(string) bool) (string, error) {
	config, err := awsConfig.LoadDefaultConfig(context.Background(), awsConfig.WithCredentialsProvider(aws.AnonymousCredentials{}))
	if err != nil {
		return "", err
	}

	s3Client := s3.NewFromConfig(config)

	// Manual URL example: https://s3.amazonaws.com/dd-agent-mstesting?prefix=pipelines/A7/25309493
	result, err := s3Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(fmt.Sprintf("pipelines/A%s/%s", majorVersion, pipelineID)),
	})

	if err != nil {
		return "", err
	}

	if len(result.Contents) <= 0 {
		return "", fmt.Errorf("no artifact found for pipeline %v", pipelineID)
	}

	for _, obj := range result.Contents {
		if !predicate(*obj.Key) {
			continue
		}

		return fmt.Sprintf("https://s3.amazonaws.com/%s/%s", bucket, *obj.Key), nil
	}

	return "", fmt.Errorf("no agent artifact found for pipeline %v", pipelineID)
}
