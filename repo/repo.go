package repo

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/session"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Repo represents the shared values needed for access to the database.
type Repo struct {
	svc *dynamodb.DynamoDB
}

// New will return you a new instance of the repo object.
func New() *Repo {
	log.Println("instantiating a new repo")
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region: aws.String("us-east-2"),
		},
		SharedConfigState: session.SharedConfigEnable,
	}))

	log.Println("repo: Succesfuly created a new repo object")
	return &Repo{
		svc: dynamodb.New(sess),
	}
}
