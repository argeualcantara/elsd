package routingkeys

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"log"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

const (
	getProjectionExpression = "StackId, Tags"
	localEndpoint = "http://localhost:8000"
	region = "us-west-2"
)

// Service provides the service object
type Service struct {
	session   *session.Session
	client    *dynamodb.DynamoDB
	tableName *string
}

// Entity represents a RoutingKey record
type Entity struct {
	ID     	string  `json:"id"`
	ServiceInstances []ServiceInstance `json:"serviceInstance"`
}

type ServiceInstance struct {
	Uri string   `json:"uri"`
	Tags []string `json:"tags"`
}


func (s *Service) createTable() (*dynamodb.CreateTableOutput, error) {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(*(s.tableName)),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("Uri"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType: aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Uri"),
				KeyType: aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits: aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	return s.client.CreateTable(params)
}

// New creates a new RoutingKeysService
func New(tableName string, id string, secret string, token string) *Service {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	creds := credentials.NewStaticCredentials(id, secret, token)

	localConfig := aws.NewConfig().
		WithCredentials(creds).
		WithEndpoint(localEndpoint).
		WithRegion(region)

	svc :=  Service{
		session:   sess,
		client:    dynamodb.New(sess, localConfig),
		tableName: &tableName,
	}

	// create table in dynamo, will fail if the table is already there
	_, err = svc.createTable()
	if err !=nil {
		log.Println("Error creating table: ",  tableName, " error: ", err)
	}

	return &svc
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

	return s.fromDynamoToEntity(id, entity)
}


func (service *Service) Add(instance *ServiceInstance, i string) (*Entity, error) {
	instances := make([]ServiceInstance,1)
	instances = append(instances, ServiceInstance{instance.Uri, instance.Tags})

	entity := Entity{i,instances}

	item, err := dynamodbattribute.MarshalMap(entity)
	if err !=nil {
		log.Println("Failed to convert", err)
		return nil, err
	}

	result, err := service.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String("RoutingKeys"),
	})
	if err !=nil {
		log.Println("Failed to write item", err)
		return nil, err
	}
	return &entity, nil
}


func (s *Service) fromDynamoToEntity(id string, input *dynamodb.QueryOutput) *Entity {
	length := len(input.Items)
	if length == 0 {
		return nil
	}

	stacks := make([]ServiceInstance, (length - 1))
	for _, value := range input.Items {
		stacks = append(stacks, ServiceInstance{Uri: *value["Uri"].S, Tags: value["Tags"].SS})
	}

	return &Entity{
		ID:     id,
		ServiceInstances: stacks,
	}
}

