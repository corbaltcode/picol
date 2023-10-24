package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/corbaltcode/picol/internal/ddbutil"
)

func MaybeUpdateSequence(ctx context.Context, ddbClient *dynamodb.Client, sequenceName string, nextId int64) error {
	tablePrefix := CtxGetDynamoDBTablePrefix(ctx)
	sequenceTableName := aws.String(tablePrefix + "Sequences")

	log.Printf("MaybeUpdateSequence: sequenceTableName=%s, sequenceName=%s, nextId=%d\n", *sequenceTableName, sequenceName, nextId)

	uii := dynamodb.UpdateItemInput{
		TableName: sequenceTableName,
		Key: map[string]ddbTypes.AttributeValue{
			"SequenceName": ddbutil.S(sequenceName),
		},
		ExpressionAttributeValues: map[string]ddbTypes.AttributeValue{
			":NextId": ddbutil.N(nextId),
		},
		UpdateExpression:    aws.String("SET NextId = :NextId"),
		ConditionExpression: aws.String("attribute_not_exists(SequenceName) OR NextId < :NextId"),
	}

	_, err := ddbClient.UpdateItem(ctx, &uii)
	if err != nil {
		var ccfe *ddbTypes.ConditionalCheckFailedException
		if errors.As(err, &ccfe) {
			return nil
		}
	}

	return err
}
