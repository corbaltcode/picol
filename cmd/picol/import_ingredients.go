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

func importIngredients(ctx context.Context, args []string) int {
	flags := flag.NewFlagSet("import-ingredients", flag.ExitOnError)
	allowUpdate := flags.Bool("allow-update", false, "Allow updating existing ingredients.")
	idSequenceOnly := flags.Bool("id-sequence-only", false, "Only update the id sequence, do not import ingredients.")
	help := flags.Bool("help", false, "Show help.")

	flags.Usage = func() {
		out := flags.Output()
		fmt.Fprintf(out, "Import ingredient data from a JSON file.\n")
		fmt.Fprintf(out, "Usage: %s import-ingredients [options] <filename>\n", os.Args[0])
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
	var ingredients picolApiV1.Response[picolApiV1.Ingredient]
	err := decoder.Decode(&ingredients)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding JSON: %s\n", err)
		return 1
	}

	awsConfig := CtxGetAWSConfig(ctx)
	tablePrefix := CtxGetDynamoDBTablePrefix(ctx)
	ddbClient := dynamodb.NewFromConfig(awsConfig)
	highestIngredientId := 0
	ingredientsUII := dynamodb.UpdateItemInput{
		TableName: aws.String(fmt.Sprintf("%sIngredients", tablePrefix)),
		ExpressionAttributeNames: map[string]string{
			"#ResistanceId":   "ResistanceId",
			"#Name":           "Name",
			"#Code":           "Code",
			"#Notes":          "Notes",
			"#ManagementCode": "ManagementCode",
		},
	}

	resistancesUII := dynamodb.UpdateItemInput{
		TableName: aws.String(fmt.Sprintf("%sResistances", tablePrefix)),
		ExpressionAttributeNames: map[string]string{
			"#Ingredients": "Ingredients",
		},
		UpdateExpression: aws.String("ADD #Ingredients :Ingredient"),
	}

	setAll := aws.String("SET #ResistanceId = :ResistanceId, #Name = :Name, #Code = :Code, #Notes = :Notes REMOVE #ManagementCode")
	setAllRemoveNotes := aws.String("SET #ResistanceId = :ResistanceId, #Name = :Name, #Code = :Code REMOVE #ManagementCode, #Notes")

	if !*allowUpdate {
		ingredientsUII.ConditionExpression = aws.String("attribute_not_exists(Id)")
	}

	for _, apiIngredient := range ingredients.Data {
		fmt.Printf("%#v\n", apiIngredient)

		if apiIngredient.Id > highestIngredientId {
			highestIngredientId = apiIngredient.Id
		}

		if *idSequenceOnly {
			continue
		}

		ingredientsUII.Key = map[string]ddbTypes.AttributeValue{
			"Id": ddbutil.N(int64(apiIngredient.Id)),
		}

		ingredientsUII.ExpressionAttributeValues = map[string]ddbTypes.AttributeValue{
			":ResistanceId": ddbutil.N(int64(apiIngredient.Resistance.Id)),
			":Name":         ddbutil.S(apiIngredient.Name),
			":Code":         ddbutil.S(apiIngredient.Code),
		}

		if apiIngredient.Notes == "" {
			ingredientsUII.ExpressionAttributeValues[":Notes"] = ddbutil.S(apiIngredient.Notes)
			ingredientsUII.UpdateExpression = setAll
		} else {
			ingredientsUII.UpdateExpression = setAllRemoveNotes
		}

		_, err = ddbClient.UpdateItem(ctx, &ingredientsUII)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error writing ingredient: %s\n", err)
			return 1
		}

		if apiIngredient.Resistance.Code != "" { // Don't update the null item
			resistancesUII.Key = map[string]ddbTypes.AttributeValue{
				"Id": ddbutil.N(int64(apiIngredient.Resistance.Id)),
			}

			resistancesUII.ExpressionAttributeValues = map[string]ddbTypes.AttributeValue{
				":Ingredient": ddbutil.NS1(int64(apiIngredient.Id)),
			}

			_, err = ddbClient.UpdateItem(ctx, &resistancesUII)

			if err != nil {
				fmt.Fprintf(os.Stderr, "Error writing ingredient to resistances: %s\n", err)
				return 1
			}
		}
	}

	sequenceName := fmt.Sprintf("%sIngredients.Id", tablePrefix)
	err = MaybeUpdateSequence(ctx, ddbClient, sequenceName, int64(highestIngredientId)+1)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating sequence: %s\n", err)
		return 1
	}

	return 0
}
