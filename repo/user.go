package repo

import (
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/uuid"
)

type User struct {
	Email    string `dynamodbav:"email"`
	Password string `dynamodbav:"password"`
	Session  string `dynamodbav:"session"`
}

func (r *Repo) GetUserBySession(session string) (User, error) {
	log.Printf("Getting user with session %s. \n", session)
	resp, err := r.svc.Query(&dynamodb.QueryInput{
		TableName: aws.String("millennium-falcon-auction-users"),
		IndexName: aws.String("session-index"),
		KeyConditions: map[string]*dynamodb.Condition{
			"session": {
				ComparisonOperator: aws.String("EQ"),
				AttributeValueList: []*dynamodb.AttributeValue{
					{
						S: aws.String(session),
					},
				},
			},
		},
		Limit: aws.Int64(1),
	})
	if err != nil {
		log.Printf("repo: Error trying to get user by session %v \n", err)
		return User{}, err
	}
	if resp.Count == nil {
		log.Printf("repo: count from query is empty \n")
		return User{}, errors.New("could not find user with session")
	}

	log.Println("Successfully retrieved user from dynamo.")

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(resp.Items[0], &user); err != nil {
		return User{}, err
	}
	return user, nil
}

func (r *Repo) GetUser(email string) (User, error) {
	log.Printf("Getting user with email %s. \n", email)
	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("millennium-falcon-auction-users"),
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
	})
	if err != nil {
		fmt.Printf("repo: Error getting user %v \n", err)
		return User{}, errors.New("could not retrieve user from dynamo")
	}

	log.Println("Successfully retrieved user from dynamo.")

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &user); err != nil {
		fmt.Printf("repo: Error marshaling into user object %v \n", err)
		return User{}, err
	}
	return user, nil
}

func (r *Repo) CreateUser(email, password string) (User, error) {
	log.Printf("repo: Attempting to create user with email %s \n", email)
	user := User{
		Email:    email,
		Password: password,
		Session:  uuid.New().String(),
	}

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		return User{}, errors.New("could not marshal created user into dynamo map")
	}

	if _, err := r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("millennium-falcon-auction-users"),
		Item:      item,
	}); err != nil {
		return User{}, errors.New("could not put created user into dynamo")
	}

	log.Printf("repo: Successfully created user %s \n", email)

	return user, nil
}
