package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/logging"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func main() {
	cliFlags := flag.NewFlagSet("picol", flag.ExitOnError)
	help := cliFlags.Bool("help", false, "Show help.")
	tablePrefix := cliFlags.String("table-prefix", "", "DynamoDB table prefix to use in addition to the project and environment names.")
	project := cliFlags.String("project", "PICOL", "Project name to prefix to DynamoDB table names.")
	environment := cliFlags.String("environment", "", "Environment name to prefix to DynamoDB table names. Will be obtained from the PICOL_ENV environment variable if not specified.")
	profile := cliFlags.String("profile", "", "The AWS profile to use. Defaults to the AWS_PROFILE environment variable if not specified.")
	region := cliFlags.String("region", "", "The AWS region to use. Defaults to the AWS_REGION/AWS_DEFAULT_REGION environment variable if not specified.")
	debug := cliFlags.Bool("debug", false, "Enable debug logging.")

	cliFlags.Usage = func() {
		out := cliFlags.Output()
		fmt.Fprintf(out, "PICOL administrative command line interface.\n")
		fmt.Fprintf(out, "Usage: %s [options] <subcommand> [subcommand options] <subcommand arguments>\n", os.Args[0])
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, "Options:\n")
		cliFlags.PrintDefaults()
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, "Subcommands:\n")
		fmt.Fprintf(out, "  import-crops  Import crop data from a JSON file.\n")
	}

	cliFlags.Parse(os.Args[1:])
	if *help {
		cliFlags.SetOutput(os.Stdout)
		cliFlags.Usage()
		os.Exit(0)
	}

	if *environment == "" {
		picol_env, found := os.LookupEnv("PICOL_ENV")
		if !found {
			picol_env = "dev"
		}

		picol_env = cases.Title(language.English).String(picol_env)
		environment = &picol_env
	}

	var configOpts []func(*config.LoadOptions) error

	if *region != "" {
		configOpts = append(configOpts, config.WithRegion(*region))
	} else {
		region, found := os.LookupEnv("AWS_REGION")
		if found {
			configOpts = append(configOpts, config.WithRegion(region))
		} else {
			region, found = os.LookupEnv("AWS_DEFAULT_REGION")
			if found {
				configOpts = append(configOpts, config.WithRegion(region))
			}
		}
	}

	if *profile != "" {
		configOpts = append(configOpts, config.WithSharedConfigProfile(*profile))
	}

	if *debug {
		log.SetOutput(os.Stderr)
		log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile | log.LUTC)
		configOpts = append(configOpts, config.WithLogger(logging.StandardLogger{Logger: log.Default()}))
	}

	awsConfig, err := config.LoadDefaultConfig(context.Background(), configOpts...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading AWS configuration: %s\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, PicolCtxDynamoDBTablePrefix, fmt.Sprintf("%s%s%s", *tablePrefix, *project, *environment))
	ctx = context.WithValue(ctx, PicolCtxAWSConfig, awsConfig)

	if cliFlags.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No subcommand specified.\n")
		cliFlags.Usage()
		os.Exit(1)
	}

	switch cliFlags.Arg(0) {
	case "import-crops":
		os.Exit(importCrops(ctx, cliFlags.Args()[1:]))
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", cliFlags.Arg(0))
		cliFlags.Usage()
		os.Exit(1)
	}
}
