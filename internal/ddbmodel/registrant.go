package ddbmodel

type Registrant struct {
	Id   int
	Name string
	Url  string `dynamodbav:",omitempty"`
}
