package ddbmodel

type Sequence struct {
	TableName string
	NextId    int `dynamodbav:",omitempty"`
}
