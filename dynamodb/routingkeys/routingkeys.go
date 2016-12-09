package routingkeys

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

const (
	getProjectionExpression = "StackId, Tags"
)

// Service provides the service object
type Service struct {
	session   *session.Session
	client    *dynamodb.DynamoDB
	tableName *string
}

// Entity represents a RoutingKey record
type Entity struct {
	ID     string  `json:"id"`
	Stacks []stack `json:"stacks"`
}

type stack struct {
	Name *string   `json:"name"`
	Tags []*string `json:"tags"`
}

// New creates a new RoutingKeysService
func New(tableName string) *Service {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	return &Service{
		session:   sess,
		client:    dynamodb.New(sess, aws.NewConfig()),
		tableName: &tableName,
	}
}

// Get returns a Routing Key
func (s *Service) Get(id string) *Entity {
	params := &dynamodb.QueryInput{
		TableName:            s.tableName,
		ProjectionExpression: aws.String(getProjectionExpression),
		KeyConditions: map[string]*dynamodb.Condition{
			"Id": {
				ComparisonOperator: aws.String(dynamodb.ComparisonOperatorEq),
				AttributeValueList: []*dynamodb.AttributeValue{
					{S: aws.String(id)},
				},
			},
		},
	}

	entity, err := s.client.Query(params)
	if err != nil {
		fmt.Printf("error querying dynamodb: %s", err)
		return nil
	}

	return fromDynamoToEntity(id, entity)
}

func fromDynamoToEntity(id string, input *dynamodb.QueryOutput) *Entity {
	length := len(input.Items)
	if length == 0 {
		return nil
	}

	stacks := make([]stack, (length - 1))
	for _, value := range input.Items {
		stacks = append(stacks, stack{Name: value["StackId"].S, Tags: value["Tags"].SS})
	}

	return &Entity{
		ID:     id,
		Stacks: stacks,
	}
}
