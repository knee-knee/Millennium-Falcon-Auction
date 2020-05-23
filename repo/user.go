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

const usersTableID = "millennium-falcon-auction-users"

// GetUserBySession will return a user based off their active session.
func (r *Repo) GetUserBySession(session string) (User, error) {
	log.Printf("repo: Getting user with session %s. \n", session)
	resp, err := r.svc.Query(&dynamodb.QueryInput{
		TableName: aws.String(usersTableID),
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
	if *resp.Count == 0 {
		log.Printf("repo: No user exists with session %s \n", session)
		return User{}, errors.New("no user exists with the session provided")
	}

	log.Println("repo: Successfully retrieved user from dynamo.")

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(resp.Items[0], &user); err != nil {
		log.Printf("repo: Error trying to unmarshal dyanmo output %v \n", err)
		return User{}, err
	}
	return user, nil
}

// GetUser will return a user object based off of their email.
func (r *Repo) GetUser(email string) (User, error) {
	log.Printf("repo: Getting user with email %s. \n", email)
	queryOutput, err := r.svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(usersTableID),
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

	log.Println("repo: Successfully retrieved user from dynamo.")

	user := User{}
	if err := dynamodbattribute.UnmarshalMap(queryOutput.Item, &user); err != nil {
		fmt.Printf("repo: Error marshaling into user object %v \n", err)
		return User{}, err
	}
	return user, nil
}

// CreateUser will create a new user in dyanmo.
func (r *Repo) CreateUser(email, password string) (User, error) {
	log.Printf("repo: Attempting to create user with email %s \n", email)
	user := User{
		Email:    email,
		Password: password,
		Session:  uuid.New().String(),
	}

	item, err := dynamodbattribute.MarshalMap(user)
	if err != nil {
		log.Printf("repo: Error marshaling user %v \n", err)
		return User{}, errors.New("could not marshal created user into dynamo map")
	}

	if _, err := r.svc.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(usersTableID),
		Item:      item,
	}); err != nil {
		log.Printf("repo: Error creating user %v \n", err)
		return User{}, errors.New("could not put created user into dynamo")
	}

	log.Printf("repo: Successfully created user %s \n", email)

	return user, nil
}
