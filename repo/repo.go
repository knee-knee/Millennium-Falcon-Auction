package repo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type Repo struct {
	svc *dynamodb.DynamoDB
}

func New() *Repo {
	log.Println("instantiating a new repo")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-2"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &Repo{
		svc: dynamodb.New(sess),
	}
}
