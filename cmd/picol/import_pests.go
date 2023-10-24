package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	picolApiV1 "github.com/corbaltcode/picol/internal/api_model/v1"
	"github.com/corbaltcode/picol/internal/ddbutil"
)

func importPests(ctx context.Context, args []string) int {
	flags := flag.NewFlagSet("import-pests", flag.ExitOnError)
	allowUpdate := flags.Bool("allow-update", false, "Allow updating existing pests.")
	idSequenceOnly := flags.Bool("id-sequence-only", false, "Only update the id sequence, do not import pests.")
	help := flags.Bool("help", false, "Show help.")

	flags.Usage = func() {
		out := flags.Output()
		fmt.Fprintf(out, "Import pest data from a JSON file.\n")
		fmt.Fprintf(out, "Usage: %s import-pests [options] <filename>\n", os.Args[0])
		fmt.Fprintf(out, "\n")
		fmt.Fprintf(out, "Options:\n")
		flags.PrintDefaults()
	}

	flags.Parse(args)

	if *help {
		flags.SetOutput(os.Stdout)
		flags.Usage()
		return 0
	}

	args = flags.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No filename specified.\n")
		flags.Usage()
		return 1
	}

	if len(args) > 1 {
		fmt.Fprintf(os.Stderr, "Unknown argument: %s\n", args[1])
		flags.Usage()
		return 1
	}

	filename := args[0]
	var input io.Reader

	if filename == "-" {
		// Read from stdin
		input = os.Stdin
	} else {
		// Read from file
		fd, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			return 1
		}
		defer fd.Close()
		input = fd
	}

	decoder := json.NewDecoder(input)
	var pests picolApiV1.Response[picolApiV1.Pest]
	err := decoder.Decode(&pests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %s\n", err)
		return 1
	}

	awsConfig := CtxGetAWSConfig(ctx)
	tablePrefix := CtxGetDynamoDBTablePrefix(ctx)
	ddbClient := dynamodb.NewFromConfig(awsConfig)
	highestPestId := 0
	uii := dynamodb.UpdateItemInput{
		TableName: aws.String(fmt.Sprintf("%sPests", tablePrefix)),
		ExpressionAttributeNames: map[string]string{
			"#Name":  "Name",
			"#Code":  "Code",
			"#Notes": "Notes",
		},
	}

	setAll := aws.String("SET #Name = :Name, #Code = :Code, #Notes = :Notes")
	setAllRemoveNotes := aws.String("SET #Name = :Name, #Code = :Code REMOVE #Notes")

	if !*allowUpdate {
		uii.ConditionExpression = aws.String("attribute_not_exists(Id)")
	}

	for _, apiPest := range pests.Data {
		fmt.Printf("%#v\n", apiPest)

		if apiPest.Id > highestPestId {
			highestPestId = apiPest.Id
		}

		if *idSequenceOnly {
			continue
		}

		uii.Key = map[string]ddbTypes.AttributeValue{
			"Id": ddbutil.N(int64(apiPest.Id)),
		}

		uii.ExpressionAttributeValues = map[string]ddbTypes.AttributeValue{
			":Name": ddbutil.S(apiPest.Name),
			":Code": ddbutil.S(apiPest.Code),
		}

		if apiPest.Notes == "" {
			uii.ExpressionAttributeValues[":Notes"] = ddbutil.S(apiPest.Notes)
			uii.UpdateExpression = setAll
		} else {
			uii.UpdateExpression = setAllRemoveNotes
		}

		_, err = ddbClient.UpdateItem(ctx, &uii)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing pest: %s\n", err)
			return 1
		}
	}

	sequenceName := fmt.Sprintf("%spests.Id", tablePrefix)
	err = MaybeUpdateSequence(ctx, ddbClient, sequenceName, int64(highestPestId)+1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating sequence: %s\n", err)
		return 1
	}

	return 0
}
