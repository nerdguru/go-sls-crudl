package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

// ItemInfo has more data for our movie item
type ItemInfo struct {
	Plot   string  `json:"plot"`
	Rating float64 `json:"rating"`
}

// Item has fields for the DynamoDB keys (Year and Title) and an ItemInfo for more data
type Item struct {
	Year  int      `json:"year"`
	Title string   `json:"title"`
	Info  ItemInfo `json:"info"`
}

// GetByYearTitle wraps up the DynamoDB calls to fetch a specific Item
// Based on https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/read_item.go
func GetByYearTitle(year, title string) (Item, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	item := Item{}

	// Perform the query
	fmt.Println("Trying to read from table: ", os.Getenv("TABLE_NAME"))
	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"year": {
				N: aws.String(year),
			},
			"title": {
				S: aws.String(title),
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	// Unmarshall the result in to an Item
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		fmt.Println(err.Error())
		return item, err
	}

	return item, nil
}

// ListByYear wraps up the DynamoDB calls to list all items of a particular year
// Based on https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/scan_items.go
func ListByYear(year string) ([]Item, error) {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)
	items := []Item{}

	// Create the Expression to fill the input struct with.
	yearAsInt, err := strconv.Atoi(year)
	filt := expression.Name("year").Equal(expression.Value(yearAsInt))

	// Get back the title, year, and rating
	proj := expression.NamesList(expression.Name("title"), expression.Name("year"))

	expr, err := expression.NewBuilder().WithFilter(filt).WithProjection(proj).Build()

	if err != nil {
		fmt.Println("Got error building expression:")
		fmt.Println(err.Error())
		return items, err
	}

	// Build the query input parameters
	params := &dynamodb.ScanInput{
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
		TableName:                 aws.String(os.Getenv("TABLE_NAME")),
	}

	// Make the DynamoDB Query API call
	result, err := svc.Scan(params)
	fmt.Println("Result", result)

	if err != nil {
		fmt.Println("Query API call failed:")
		fmt.Println((err.Error()))
		return items, err
	}

	numItems := 0
	for _, i := range result.Items {
		item := Item{}

		err = dynamodbattribute.UnmarshalMap(i, &item)

		if err != nil {
			fmt.Println("Got error unmarshalling:")
			fmt.Println(err.Error())
			return items, err
		}

		fmt.Println("Title: ", item.Title)
		items = append(items, item)
		numItems++
	}

	fmt.Println("Found", numItems, "movie(s) in year ", year)
	if err != nil {
		fmt.Println(err.Error())
		return items, err
	}

	return items, nil
}

// Post extracts the Item JSON and writes it to DynamoDB
// Based on https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/create_item.go
func Post(body string) (Item, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var thisItem Item
	json.Unmarshal([]byte(body), &thisItem)

	// Take out non-alphanumberic except space characters from the title for easier slug building on reads
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	thisItem.Title = reg.ReplaceAllString(thisItem.Title, "")
	fmt.Println("Item to add:", thisItem)

	// Marshall the Item into a Map DynamoDB can deal with
	av, err := dynamodbattribute.MarshalMap(thisItem)
	if err != nil {
		fmt.Println("Got error marshalling map:")
		fmt.Println(err.Error())
		return thisItem, err
	}

	// Create Item in table and return
	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}
	_, err = svc.PutItem(input)
	return thisItem, err

}

// Delete wraps up the DynamoDB calls to delete a specific Item
// Based on https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/delete_item.go
func Delete(year, title string) error {
	// Build the Dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Perform the delete
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"year": {
				N: aws.String(year),
			},
			"title": {
				S: aws.String(title),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err := svc.DeleteItem(input)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

// Put extracts the Item JSON and updates it in DynamoDB
// Based on https://github.com/awsdocs/aws-doc-sdk-examples/blob/master/go/example_code/dynamodb/update_item.go
func Put(body string) (Item, error) {
	// Create the dynamo client object
	sess := session.Must(session.NewSession())
	svc := dynamodb.New(sess)

	// Marshall the requrest body
	var thisItem Item
	json.Unmarshal([]byte(body), &thisItem)

	// Take out non-alphanumberic except space characters from the title for easier slug building on reads
	reg, err := regexp.Compile("[^a-zA-Z0-9\\s]+")
	thisItem.Title = reg.ReplaceAllString(thisItem.Title, "")
	fmt.Println("Item to update:", thisItem)

	// Update Item in table and return
	input := &dynamodb.UpdateItemInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":r": {
				N: aws.String(strconv.FormatFloat(thisItem.Info.Rating, 'f', 1, 64)),
			},
			":p": {
				S: aws.String(thisItem.Info.Plot),
			},
		},
		TableName: aws.String(os.Getenv("TABLE_NAME")),
		Key: map[string]*dynamodb.AttributeValue{
			"year": {
				N: aws.String(strconv.Itoa(thisItem.Year)),
			},
			"title": {
				S: aws.String(thisItem.Title),
			},
		},
		ReturnValues:     aws.String("UPDATED_NEW"),
		UpdateExpression: aws.String("set info.rating = :r, info.plot = :p"),
	}

	_, err = svc.UpdateItem(input)
	if err != nil {
		fmt.Println(err.Error())
	}
	return thisItem, err

}
