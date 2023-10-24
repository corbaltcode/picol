package ddbmodel

type Crop struct {
	Id    int
	Code  string
	Name  string
	Notes string `dynamodbav:",omitempty"`
}
