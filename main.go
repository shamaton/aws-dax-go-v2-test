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
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type User struct {
	UserId    string                 `dynamodbav:"user_id"`
	GameTitle string                 `dynamodbav:"game_title"`
	Score     int                    `dynamodbav:"score"`
	Info      map[string]interface{} `dynamodbav:"info"`
}

func (u User) GetKey() map[string]types.AttributeValue {
	av, err := attributevalue.MarshalMap(u)
	if err != nil {
		panic(fmt.Errorf("attribute value marshal error. user: %s, title: %s, err: %v", u.UserId, u.GameTitle, err))
	}
	delete(av, "score")
	delete(av, "info")
	return av
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

	written, err := daxCli.BatchWriteItems(ctx, []User{
		{UserId: "user-a", GameTitle: "gt-a", Score: 1, Info: map[string]interface{}{"rating": 1.1, "comment": "aaa"}},
		{UserId: "user-b", GameTitle: "gt-b", Score: 2, Info: map[string]interface{}{"rating": 2.2, "comment": "bbb"}},
		{UserId: "user-c", GameTitle: "gt-c", Score: 3, Info: map[string]interface{}{"rating": 3.3, "comment": "ccc"}},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] batch write users count is", written)

	users, _, err := daxCli.BatchGetItems(ctx, []User{
		{UserId: "user-b", GameTitle: "gt-b"},
		{UserId: "user-c", GameTitle: "gt-c"},
		{UserId: "user-d", GameTitle: "gt-d"},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] batch get users is", users)

	users, err = daxCli.Query(ctx, "gt-c")
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] query users is", users)

	users, err = daxCli.Scan(ctx, 1, 2)
	if err != nil {
		panic(err)
	}
	fmt.Println("[dax] scan users is", users)

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

	in := &dynamodb.GetItemInput{
		TableName: &c.tableName,
		Key:       user.GetKey(),
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
		Info:      map[string]interface{}{},
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

	in := &dynamodb.DeleteItemInput{
		TableName: &c.tableName,
		Key:       user.GetKey(),
	}
	if len(cb) > 0 {
		cb[0](in)
	}
	_, err := c.c.DeleteItem(ctx, in)
	if err != nil {
		return fmt.Errorf("delete item error: %v", err)
	}
	return nil
}

func (c *client) UpdateItem(ctx context.Context, user User) (*User, error) {
	update := expression.Set(expression.Name("info.rating"), expression.Value(user.Info["rating"]))
	update.Set(expression.Name("info.comment"), expression.Value(user.Info["comment"]))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expression for update. reason: %v", err)
	}

	in := &dynamodb.UpdateItemInput{
		TableName:                 &c.tableName,
		Key:                       user.GetKey(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueUpdatedNew,
	}

	out, err := c.c.UpdateItem(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("couldn't update user %v. reason: %v", user.UserId, err)
	}

	user = User{}
	err = attributevalue.UnmarshalMap(out.Attributes, &user)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshall update out. reason: %v", err)
	}

	return &user, err
}

func (c *client) BatchGetItems(ctx context.Context, users []User) ([]User, int, error) {
	var (
		gotten    = 0
		batchSize = 100 // DynamoDB allows a maximum batch size of 100 items.
		start     = 0
		end       = start + batchSize
	)

	results := make([]User, 0, len(users))
	for start < len(users) {
		if end > len(users) {
			end = len(users)
		}

		items := types.KeysAndAttributes{
			Keys: make([]map[string]types.AttributeValue, 0, end-start),
		}
		for _, user := range users[start:end] {
			items.Keys = append(items.Keys, user.GetKey())
		}

		in := &dynamodb.BatchGetItemInput{
			RequestItems: map[string]types.KeysAndAttributes{c.tableName: items},
		}
		out, err := c.c.BatchGetItem(ctx, in)
		if err != nil {
			return nil, gotten, fmt.Errorf("couldn't get a batch of users to %v. reason: %v", c.tableName, err)
		}

		resps, ok := out.Responses[c.tableName]
		if !ok {
			return nil, gotten, fmt.Errorf("not found response to %s", c.tableName)
		}

		// should retry unprocessed keys
		if len(out.UnprocessedKeys) > 0 {
			return nil, gotten, fmt.Errorf("found unprocessed keys. keys: %v", out.UnprocessedKeys)
		}

		var gotUsers []User
		if err = attributevalue.UnmarshalListOfMaps(resps, &gotUsers); err != nil {
			return nil, gotten, fmt.Errorf("couldn't unmarshal query out. reason: %v", err)
		}

		gotten += len(gotUsers)
		start = end
		end += batchSize
		results = append(results, gotUsers...)
	}

	return results, gotten, nil
}

func (c *client) BatchWriteItems(ctx context.Context, users []User) (int, error) {
	var (
		written   = 0
		batchSize = 25 // DynamoDB allows a maximum batch size of 25 items.
		start     = 0
		end       = start + batchSize
	)
	for start < len(users) {
		if end > len(users) {
			end = len(users)
		}

		reqs := make([]types.WriteRequest, 0, end-start)
		for _, user := range users[start:end] {
			item, err := attributevalue.MarshalMap(user)
			if err != nil {
				return written, fmt.Errorf("couldn't marshal user %v for batch writing. reason: %v", user.UserId, err)
			}
			req := types.WriteRequest{PutRequest: &types.PutRequest{Item: item}}
			reqs = append(reqs, req)
		}

		in := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{c.tableName: reqs},
		}
		_, err := c.c.BatchWriteItem(ctx, in)
		if err != nil {
			return written, fmt.Errorf("couldn't add a batch of users to %v. reason: %v", c.tableName, err)
		}

		written += len(reqs)
		start = end
		end += batchSize
	}

	return written, nil
}

func (c *client) Query(ctx context.Context, userId string, cb ...func(*dynamodb.QueryInput)) ([]User, error) {
	var user []User
	keyEx := expression.Key("user_id").Equal(expression.Value(userId))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expression for query. reason: %v", err)
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
		return nil, fmt.Errorf("couldn't query for user in %v. reason: %v", userId, err)
	}

	if err = attributevalue.UnmarshalListOfMaps(out.Items, &user); err != nil {
		return nil, fmt.Errorf("couldn't unmarshal query out. reason: %v", err)
	}

	return user, err
}

func (c *client) Scan(ctx context.Context, startScore, endScore int, cb ...func(*dynamodb.ScanInput)) ([]User, error) {
	var users []User

	filter := expression.Name("score").Between(expression.Value(startScore), expression.Value(endScore))
	projection := expression.NamesList(
		expression.Name("user_id"),
		expression.Name("game_title"),
		expression.Name("info.rating"),
	)
	expr, err := expression.NewBuilder().WithFilter(filter).WithProjection(projection).Build()
	if err != nil {
		return nil, fmt.Errorf("couldn't build expressions for scan. reason: %v", err)
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
		return nil, fmt.Errorf("couldn't scan for users released between %v and %v. reason: %v",
			startScore, endScore, err)
	}

	err = attributevalue.UnmarshalListOfMaps(out.Items, &users)
	if err != nil {
		log.Printf("Couldn't unmarshal query out. reason: %v\n", err)
	}
	return users, err
}
