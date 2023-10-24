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

func importResistances(ctx context.Context, args []string) int {
	flags := flag.NewFlagSet("import-resistances", flag.ExitOnError)
	allowUpdate := flags.Bool("allow-update", false, "Allow updating existing resistances.")
	idSequenceOnly := flags.Bool("id-sequence-only", false, "Only update the id sequence, do not import resistances.")
	help := flags.Bool("help", false, "Show help.")
	clearIngredients := flags.Bool("clear-ingredients", true, "Clear the ingredients list for each imported resistance.")

	flags.Usage = func() {
		out := flags.Output()
		fmt.Fprintf(out, "Import resistance data from a JSON file.\n")
		fmt.Fprintf(out, "Usage: %s import-resistances [options] <filename>\n", os.Args[0])
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
	var resistances picolApiV1.Response[picolApiV1.Resistance]
	err := decoder.Decode(&resistances)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %s\n", err)
		return 1
	}

	awsConfig := CtxGetAWSConfig(ctx)
	tablePrefix := CtxGetDynamoDBTablePrefix(ctx)

	ddbClient := dynamodb.NewFromConfig(awsConfig)
	tableName := aws.String(fmt.Sprintf("%sResistances", tablePrefix))

	highestResistanceId := 0
	uii := dynamodb.UpdateItemInput{
		TableName: tableName,
		ExpressionAttributeNames: map[string]string{
			"#Source":         "Source",
			"#Code":           "Code",
			"#MethodOfAction": "MethodOfAction",
		},
	}

	if *clearIngredients {
		uii.UpdateExpression = aws.String("SET #Source = :Source, #Code = :Code, #MethodOfAction = :MethodOfAction REMOVE #Ingredients")
		uii.ExpressionAttributeNames["#Ingredients"] = "Ingredients"
	} else {
		uii.UpdateExpression = aws.String("SET #Source = :Source, #Code = :Code, #MethodOfAction = :MethodOfAction")
	}

	if !*allowUpdate {
		uii.ConditionExpression = aws.String("attribute_not_exists(Id)")
	}

	for _, apiResistance := range resistances.Data {
		fmt.Printf("%#v\n", apiResistance)

		if apiResistance.Id > highestResistanceId {
			highestResistanceId = apiResistance.Id
		}

		if *idSequenceOnly {
			continue
		}

		uii.Key = map[string]ddbTypes.AttributeValue{
			"Id": ddbutil.N(int64(apiResistance.Id)),
		}

		uii.ExpressionAttributeValues = map[string]ddbTypes.AttributeValue{
			":Source":         ddbutil.S(apiResistance.Source),
			":Code":           ddbutil.S(apiResistance.Code),
			":MethodOfAction": ddbutil.S(apiResistance.MethodOfAction),
		}

		_, err = ddbClient.UpdateItem(ctx, &uii)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing ingredient: %s\n", err)
			return 1
		}
	}

	sequenceName := fmt.Sprintf("%sResistances.Id", tablePrefix)
	err = MaybeUpdateSequence(ctx, ddbClient, sequenceName, int64(highestResistanceId)+1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating sequence: %s\n", err)
		return 1
	}

	return 0
}
