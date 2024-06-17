package pinpoint

import (
	"context"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/pinpointsmsvoicev2"
	"github.com/breathbath/goalert/config"
	"os"
)

func NewClient(
	ctx context.Context,
	cfg config.Config,
	opts *Config,
) (*pinpointsmsvoicev2.Client, error) {
	if cfg.PinPoint.AwsAccessKeyId != "" {
		err := os.Setenv("AWS_ACCESS_KEY_ID", cfg.PinPoint.AwsAccessKeyId)
		if err != nil {
			return nil, err
		}
	}

	if cfg.PinPoint.AwsSecretAccessKey != "" {
		err := os.Setenv("AWS_SECRET_ACCESS_KEY", cfg.PinPoint.AwsSecretAccessKey)
		if err != nil {
			return nil, err
		}
	}

	if cfg.PinPoint.AwsSessionToken != "" {
		err := os.Setenv("AWS_SESSION_TOKEN", cfg.PinPoint.AwsSessionToken)
		if err != nil {
			return nil, err
		}
	}

	optFns := make([]func(options *awsConfig.LoadOptions) error, 0)
	if opts.Client != nil {
		optFns = append(optFns, awsConfig.WithHTTPClient(opts.Client))
	}

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return nil, err
	}

	if opts.BaseURL != "" {
		awsCfg.BaseEndpoint = &opts.BaseURL
	}

	return pinpointsmsvoicev2.NewFromConfig(awsCfg), nil
}
