package ddbutil

import (
	"strconv"

	ddbTypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// S creates a new DynamoDB string attribute value member.
func S(s string) *ddbTypes.AttributeValueMemberS {
	return &ddbTypes.AttributeValueMemberS{Value: s}
}

// L creates a new DynamoDB list attribute value member from a list of attribute value members.
func L(l []ddbTypes.AttributeValue) *ddbTypes.AttributeValueMemberL {
	return &ddbTypes.AttributeValueMemberL{Value: l}
}

// M creates a new DynamoDB map attribute value member from a map.
func M(m map[string]ddbTypes.AttributeValue) *ddbTypes.AttributeValueMemberM {
	return &ddbTypes.AttributeValueMemberM{Value: m}
}

// N creates a new DynamoDB number attribute value member from an integer.
func N(i int64) *ddbTypes.AttributeValueMemberN {
	return &ddbTypes.AttributeValueMemberN{Value: strconv.FormatInt(i, 10)}
}

// BOOL creates a new DynamoDB boolean attribute value member.
func BOOL(b bool) *ddbTypes.AttributeValueMemberBOOL {
	return &ddbTypes.AttributeValueMemberBOOL{Value: b}
}

// NULL creates a new DynamoDB null value member.
func NULL(isNull bool) *ddbTypes.AttributeValueMemberNULL {
	return &ddbTypes.AttributeValueMemberNULL{Value: isNull}
}

// NS creates a new DynamoDB number set attribute value member from a list of integers.
func NS(ns []int64) *ddbTypes.AttributeValueMemberNS {
	s := make([]string, len(ns))
	for i, n := range ns {
		s[i] = strconv.FormatInt(n, 10)
	}
	return &ddbTypes.AttributeValueMemberNS{Value: s}
}

// NS1 creates a new DynamoDB number set attribute value member from a single integer.
func NS1(n int64) *ddbTypes.AttributeValueMemberNS {
	s := make([]string, 1)
	s[0] = strconv.FormatInt(n, 10)
	return &ddbTypes.AttributeValueMemberNS{Value: s}
}
