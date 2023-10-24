package ddbmodel

type Resistance struct {
	Id             int
	Source         string
	Code           string
	MethodOfAction string
	Ingredients    []int `dynamodbav:",numberset"`

	// Rid is not accessible
}
