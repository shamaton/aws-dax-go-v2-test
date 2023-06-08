package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-dax-go/dax"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type User struct {
	UserId    string `dynamodbav:"user_id"`
	GameTitle string `dynamodbav:"game_title"`
	Score     int    `dynamodbav:"score"`
}

var (
	endpoint  string
	tableName = "GameScores"
)

func main() {
	flag.StringVar(&endpoint, "ep", "", "dax endpoint")
	flag.Parse()

	if endpoint == "" {
		flag.Usage()
		return
	}

	ctx := context.Background()
	daxCli := daxClient()
	dynamoCli := dynamoClient()

	id := "user-a"
	title := "gt-a"

	user, err := daxCli.GetItem(ctx, "user-a", "gt-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] get user is", user)

	user, err = daxCli.PutItem(ctx, id, title, 100)
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] put user is", user)

	user, err = daxCli.GetItem(ctx, "user-a", "gt-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] get user is", user)

	// dynamo
	_, err = dynamoCli.PutItem(ctx, id, title, int(time.Now().Unix()))
	if err != nil {
		panic(err)
	}
	fmt.Println("[dynamo] put user is", user)

	user, err = dynamoCli.GetItem(ctx, "user-a", "gt-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dynamo] get user is", user)

	user, err = daxCli.GetItem(ctx, "user-a", "gt-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] get user is", user)

	err = daxCli.DeleteItem(ctx, "user-a", "gt-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] delete user")

	user, err = daxCli.GetItem(ctx, "user-a", "gt-a")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] get user is", user)
}

type client struct {
	c         dax.DynamoDBAPI
	tableName string
}

func daxClient() *client {
	cfg := aws.Config{
		Region: "us-east-1",
		EndpointResolverWithOptions: aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: endpoint,
			}, nil
		}),
	}

	ctx := context.Background()
	c, err := dax.NewWithSDKConfig(ctx, cfg)
	if err != nil {
		panic(err)
	}
	return &client{c: c, tableName: tableName}
}

func dynamoClient() *client {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		panic(err)
	}
	return &client{c: dynamodb.NewFromConfig(cfg), tableName: tableName}
}

func (c *client) GetItem(ctx context.Context, id, title string, cb ...func(*dynamodb.GetItemInput)) (*User, error) {
	user := User{
		UserId:    id,
		GameTitle: title,
	}
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, fmt.Errorf("attribute value marshal error. user: %s, title: %s, err: %v", id, title, err)
	}
	delete(av, "score")

	in := &dynamodb.GetItemInput{
		TableName: &c.tableName,
		Key:       av,
	}
	if len(cb) > 0 {
		cb[0](in)
	}
	out, err := c.c.GetItem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("get item error. user: %s, title: %s, err: %v", id, title, err)
	}

	user = User{}
	err = attributevalue.UnmarshalMap(out.Item, &user)
	//p.P("get: ", out.Item)
	if err != nil {
		return nil, fmt.Errorf("attribute value unmarshal error. user: %s, title: %s, err: %v", id, title, err)
	}
	if user.UserId != id || user.GameTitle != title {
		return nil, nil
	}
	return &user, nil
}

func (c *client) PutItem(ctx context.Context, id, title string, score int, cb ...func(input *dynamodb.PutItemInput)) (*User, error) {
	user := &User{
		UserId:    id,
		GameTitle: title,
		Score:     score,
	}
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, fmt.Errorf("attribute value marshal error: %v", err)
	}

	in := &dynamodb.PutItemInput{
		TableName: &c.tableName,
		Item:      av,
	}
	if len(cb) > 0 {
		cb[0](in)
	}
	//p.P("put: ", in.Item)
	_, err = c.c.PutItem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("put item error: %v", err)
	}
	return user, nil
}

func (c *client) DeleteItem(ctx context.Context, id, title string, cb ...func(input *dynamodb.DeleteItemInput)) error {
	user := &User{
		UserId:    id,
		GameTitle: title,
	}
	av, err := attributevalue.MarshalMap(user)
	if err != nil {
		return fmt.Errorf("attribute value marshal error: %v", err)
	}
	delete(av, "score")

	in := &dynamodb.DeleteItemInput{
		TableName: &c.tableName,
		Key:       av,
	}
	if len(cb) > 0 {
		cb[0](in)
	}
	_, err = c.c.DeleteItem(ctx, in)
	if err != nil {
		return fmt.Errorf("delete item error: %v", err)
	}
	return nil
}

func (c *client) Query(ctx context.Context, score int, cb ...func(*dynamodb.QueryInput)) ([]User, error) {
	var user []User
	keyEx := expression.Key("score").Equal(expression.Value(score))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expression for query. Here's why: %v", err)
	}

	in := &dynamodb.QueryInput{
		TableName:                 &c.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	}
	if len(cb) > 0 {
		cb[0](in)
	}

	out, err := c.c.Query(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("couldn't query for user in %v. Here's why: %v", score, err)
	}

	if err = attributevalue.UnmarshalListOfMaps(out.Items, &user); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal query out. Here's why: %v", err)
	}

	return user, err
}

func (c *client) Scan(ctx context.Context, startYear int, endYear int, cb ...func(*dynamodb.ScanInput)) ([]User, error) {
	var users []User

	filter := expression.Name("score").Between(expression.Value(startYear), expression.Value(endYear))
	projection := expression.NamesList(
		expression.Name("user_id"),
		expression.Name("game_title"),
		expression.Name("info.rating"),
	)
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expressions for scan. Here's why: %v", err)
	}

	in := &dynamodb.ScanInput{
		TableName:                 &c.tableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}
	if len(cb) > 0 {
		cb[0](in)
	}

	out, err := c.c.Scan(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("couldn't scan for users released between %v and %v. Here's why: %v",
			startYear, endYear, err)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
	if err != nil {
		log.Printf("Couldn't unmarshal query out. Here's why: %v\n", err)
	}
	return users, err
}
