package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	dynaClient dynamodbiface.DynamoDBAPI
)

const tableName = "random-article-ids"

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Hello!"), nil
}

type Item struct {
	ArticleId string
}

func main() {
	region := os.Getenv("AWS_REGION")
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		return
	}
	dynaClient = dynamodb.New(awsSession)

	resp, err := http.Get("https://content.guardianapis.com/search?api-key={API_KEY}")
	if err != nil {
		log.Fatalln(err)
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	jsonStream := buf.String()

	type Result struct {
		Id   string
		Type string
	}

	type Response struct {
		Results []Result
	}

	type Capi struct {
		Response Response
	}

	capi := &Capi{}
	err = json.Unmarshal([]byte(jsonStream), capi)
	if err != nil {
		fmt.Println(err)
	}

	results := capi.Response.Results
	filtered := []Result{}
	for i := range results {
		if results[i].Type == "article" {
			filtered = append(filtered, results[i])
		}
	}

	for _, article := range filtered {
		item := Item{
			ArticleId: article.Id,
		}

		av, err := dynamodbattribute.MarshalMap(item)

		if err != nil {
			fmt.Println("Got error marshalling new item")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		input := &dynamodb.PutItemInput{
			Item:      av,
			TableName: aws.String(tableName),
		}

		_, err = dynaClient.PutItem(input)
		if err != nil {
			fmt.Println("Got error calling PutItem:")
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}

	lambda.Start(HandleRequest)
}
