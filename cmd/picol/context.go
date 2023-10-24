package main

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// PicolCtxKey is a unique type for context keys.
type PicolCtxKey int

const (
	// PicolCtxDynamoDBTablePrefix is a context key for the DynamoDB table prefix.
	PicolCtxDynamoDBTablePrefix PicolCtxKey = iota

	// PicolCtxAWSConfig is a context key for the AWS SDK configuration.
	PicolCtxAWSConfig
)

// CtxGetDynamoDBTablePrefix returns the DynamoDB table prefix from the context.
func CtxGetDynamoDBTablePrefix(ctx context.Context) string {
	prefixAny := ctx.Value(PicolCtxDynamoDBTablePrefix)
	if prefixAny == nil {
		return ""
	}

	prefix, ok := prefixAny.(string)
	if !ok {
		panic("PicolCtxDynamoDBTablePrefix is not a string")
	}

	return prefix
}

// CtxGetAWSConfig returns the AWS SDK configuration from the context.
func CtxGetAWSConfig(ctx context.Context) aws.Config {
	configAny := ctx.Value(PicolCtxAWSConfig)
	if configAny == nil {
		return aws.Config{}
	}

	config, ok := configAny.(aws.Config)
	if !ok {
		panic("PicolCtxAWSConfig is not an aws.Config")
	}

	return config
}
