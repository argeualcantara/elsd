/**
 * (C) Copyright 2012-2016 HP Development Company, L.P.
 * Confidential computer software. Valid license from HP required for possession, use or copying.
 * Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
 * Computer Software Documentation, and Technical Data for Commercial Items are licensed
 * to the U.S. Government under vendor's standard commercial license.
 */
package routingkeys

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"log"
)

const (
	getProjectionExpression = "Id, Uri, Tags"
	region                  = "us-west-2"
)

// Service provides the s object
type Service struct {
	session   *session.Session
	client    *dynamodb.DynamoDB
	tableName string
}

// Entity represents a ServiceInstance records for a given RoutingKey id.
type Entity struct {
	Id               string
	ServiceInstances []ServiceInstance
}

// ServiceInstance represents single ServiceInstnace record
type ServiceInstance struct {
	Id   string   `dynamodbav:"Id"`
	Uri  string   `dynamodbav:"Uri"`
	Tags []string `dynamodbav:"Tags" dynamodbav:",stringset" `
}

func (s *Service) createTable() (*dynamodb.CreateTableOutput, error) {
	params := &dynamodb.CreateTableInput{
		TableName: aws.String(s.tableName),
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
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Uri"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(5),
		},
	}

	return s.client.CreateTable(params)
}

// New creates a new RoutingKeysService
func New(tableName string, dynamoAddr string, id string, secret string, token string) *Service {
	sess, err := session.NewSession()
	if err != nil {
		panic(err)
	}

	creds := credentials.NewStaticCredentials(id, secret, token)

	localConfig := aws.NewConfig().
		WithCredentials(creds).
		WithEndpoint(dynamoAddr).
		WithRegion(region)

	svc := Service{
		session:   sess,
		client:    dynamodb.New(sess, localConfig),
		tableName: tableName,
	}

	// create table in dynamo, will fail if the table is already there
	_, err = svc.createTable()
	if err != nil {
		log.Println("Error creating table: ", tableName, " error: ", err)
	}

	return &svc
}

// Get returns all ServiceInstance for a given Routing Key
func (s *Service) Get(id string) *Entity {
	params := &dynamodb.QueryInput{
		TableName:            &s.tableName,
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

	items, err := s.client.Query(params)
	if err != nil {
		fmt.Printf("error querying dynamodb: %s", err)
		return nil
	}

	return s.fromDynamoToEntity(id, items)
}

func (s *Service) Add(instance *ServiceInstance) (error) {
	item, err := dynamodbattribute.MarshalMap(instance)
	if err != nil {
		log.Println("Failed to convert", err)
		return err
	}

	_, err = s.client.PutItem(&dynamodb.PutItemInput{
		Item:      item,
		TableName: aws.String( s.tableName),
	})
	if err != nil {
		log.Println("Failed to write item", err)
		return err
	}
	return nil
}

func (s *Service) fromDynamoToEntity(id string, input *dynamodb.QueryOutput) *Entity {
	length := len(input.Items)
	if length == 0 {
		return nil
	}

	serviceInstances := []ServiceInstance{}

	err := dynamodbattribute.UnmarshalListOfMaps(input.Items, &serviceInstances)
	if err != nil {
		log.Println("Failed to convert", err)
		return nil
	}

	return &Entity{
		Id:               id,
		ServiceInstances: serviceInstances,
	}
}
