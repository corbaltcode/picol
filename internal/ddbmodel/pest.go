package ddbmodel

type Pest struct {
	Id    int
	Name  string
	Code  string
	Notes string `dynamodbav:",omitempty"`
}
