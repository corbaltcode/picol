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

func importCrops(ctx context.Context, args []string) int {
	flags := flag.NewFlagSet("import-crops", flag.ExitOnError)
	allowUpdate := flags.Bool("allow-update", false, "Allow updating existing crops.")
	idSequenceOnly := flags.Bool("id-sequence-only", false, "Only update the id sequence, do not import crops.")
	help := flags.Bool("help", false, "Show help.")

	flags.Usage = func() {
		out := flags.Output()
		fmt.Fprintf(out, "Import crop data from a JSON file.\n")
		fmt.Fprintf(out, "Usage: %s import-crops [options] <filename>\n", os.Args[0])
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
	var crops picolApiV1.Response[picolApiV1.Crop]
	err := decoder.Decode(&crops)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %s\n", err)
		return 1
	}

	awsConfig := CtxGetAWSConfig(ctx)
	tablePrefix := CtxGetDynamoDBTablePrefix(ctx)

	ddbClient := dynamodb.NewFromConfig(awsConfig)
	tableName := aws.String(fmt.Sprintf("%sCrops", tablePrefix))

	highestCropId := 0
	uii := dynamodb.UpdateItemInput{
		TableName: tableName,
		ExpressionAttributeNames: map[string]string{
			"#Code":  "Code",
			"#Name":  "Name",
			"#Notes": "Notes",
		},
	}

	setAll := aws.String("SET #Code = :Code, #Name = :Name, #Notes = :Notes")
	setAllRemoveNotes := aws.String("SET #Code = :Code, #Name = :Name REMOVE #Notes")

	if !*allowUpdate {
		uii.ConditionExpression = aws.String("attribute_not_exists(Id)")
	}

	for _, apiCrop := range crops.Data {
		fmt.Printf("%#v\n", apiCrop)

		if apiCrop.Id > highestCropId {
			highestCropId = apiCrop.Id
		}

		if *idSequenceOnly {
			continue
		}

		uii.Key = map[string]ddbTypes.AttributeValue{
			"Id": ddbutil.N(int64(apiCrop.Id)),
		}

		uii.ExpressionAttributeValues = map[string]ddbTypes.AttributeValue{
			":Code": ddbutil.S(apiCrop.Code),
			":Name": ddbutil.S(apiCrop.Name),
		}

		if apiCrop.Notes == "" {
			uii.ExpressionAttributeValues[":Notes"] = ddbutil.S(apiCrop.Notes)
			uii.UpdateExpression = setAll
		} else {
			uii.UpdateExpression = setAllRemoveNotes
		}

		_, err = ddbClient.UpdateItem(ctx, &uii)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing crop: %s\n", err)
			return 1
		}
	}

	sequenceName := fmt.Sprintf("%sCrops.Id", tablePrefix)
	err = MaybeUpdateSequence(ctx, ddbClient, sequenceName, int64(highestCropId)+1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating sequence: %s\n", err)
		return 1
	}

	return 0
}
