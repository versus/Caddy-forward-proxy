package dynamo

import (
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var errInvalidCredentials = errors.New("Invalid credentials")

type CredentialsRepository struct {
	db    *dynamodb.DynamoDB
	table *string
}

func NewCredentialsRepository(dynamo *dynamodb.DynamoDB, table string) *CredentialsRepository {
	return &CredentialsRepository{
		db:    dynamo,
		table: aws.String(table),
	}
}

func (r *CredentialsRepository) GetPassword(login string) (string, error) {
	req := r.db.GetItemRequest(&dynamodb.GetItemInput{
		AttributesToGet: []string{
			"Pass",
		},
		Key: map[string]dynamodb.AttributeValue{
			"Login": {
				S: aws.String(login),
			},
		},
		TableName: r.table,
	})
	out, err := req.Send()
	if err != nil {
		return "", err
	}

	if len(out.Item) == 0 {
		return "", errInvalidCredentials
	}

	return *out.Item["Pass"].S, nil
}
