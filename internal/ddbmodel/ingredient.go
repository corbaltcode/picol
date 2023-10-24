package ddbmodel

type Ingredient struct {
	Id             int
	ResistanceId   *int `dynamodbav:",omitempty"`
	Name           string
	Code           string
	Notes          string `dynamodbav:",omitempty"`
	ManagementCode string `dynamodbav:",omitempty"`
}
