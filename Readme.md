# go-sls-crudl

This project riffs off of the [Dynamo DB Golang samples](https://github.com/awsdocs/aws-doc-sdk-examples/tree/master/go/example_code/dynamodb) and the [Serverless Framework Go example](https://serverless.com/blog/framework-example-golang-lambda-support/) to create an example of how to build a simple API Gateway -> Lambda -> DynamoDB set of methods.

## Code Organization
Note that, instead of using the `create_table.go` to set up the initial table, the resource building mechanism that Serverless provides is used.  Individual code is organized as follows:

* functions/post.do - POST method for creating a new item
* functions/get.do - GET method for reading a specific item
* functions/delete.do - DELETE method for deleting a specific item
* functions/put.do - PUT method for updating an existing item
* functions/list.do - GET method for listing all or a subset of items
* img/* - Images of DynamoDB tables to make this Readme easier to follow
* moviedao/moviedao.do - DAO wrapper around the DynamoDB calls
* data/XXX.json - Set of sample data files for POST and PUT actions
* Makefile - Used for dep package management and compiles of individual functions
* serverless.yml - Defines the initial table, function defs, and API Gateway events

Note that given the recency of Go support on both AWS Lambda and the Serverless Framework, combined with my own Go noob-ness, I'm not entierly certain this is the best layout but it was functional.  My hope is that it helps spark a healthy debate over what a Go Serverless project should look like.

## Set Up
If you are a Serverless Framework rookie, [follow the installation instructions here](https://serverless.com/blog/anatomy-of-a-serverless-app/#setup).  If you are a grizzled vet, be sure that you have v1.26 or later as that's the version that introduces Go support.  You'll also need to [install Go](https://golang.org/doc/install).

When both of those tasks are done, cd into your `GOPATH` and clone this project into that folder.  Then cd into the resulting `go-sls-crudle` folder and compile the source with `make`:

```bash
$ make
dep ensure
env GOOS=linux go build -ldflags="-s -w" -o bin/get functions/get.go
env GOOS=linux go build -ldflags="-s -w" -o bin/post functions/post.go
env GOOS=linux go build -ldflags="-s -w" -o bin/delete functions/delete.go
env GOOS=linux go build -ldflags="-s -w" -o bin/put functions/put.go
env GOOS=linux go build -ldflags="-s -w" -o bin/list-by-year functions/list-by-year.go
```

Finally, deploy with the 'sls' command:

```bash
$ sls deploy
Serverless: Packaging service...
Serverless: Excluding development dependencies...
Serverless: Uploading CloudFormation file to S3...
Serverless: Uploading artifacts...
Serverless: Uploading service .zip file to S3 (15.19 MB)...
Serverless: Validating template...
Serverless: Updating Stack...
Serverless: Checking Stack update progress...
..................
Serverless: Stack update finished...
Service Information
service: go-crud
stage: dev
region: us-east-1
stack: go-crud-dev
api keys:
  None
endpoints:
  GET - https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/dev/go-sls-crudl/{year}
  GET - https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/dev/go-sls-crudl/{year}/{title}
  POST - https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/dev/go-sls-crudl
  DELETE - https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/dev/go-sls-crudl/{year}/{title}
  PUT - https://xxxxxxxxxx.execute-api.us-east-1.amazonaws.com/dev/go-sls-crudl
functions:
  list: go-crud-dev-list
  get: go-crud-dev-get
  post: go-crud-dev-post
  delete: go-crud-dev-delete
  put: go-crud-dev-put
Serverless: Removing old service versions...
```

When done, you can find the new DynamoDB table in the AWS Console, which should initially look like this:

![Initial DynamoDB Table](/img/initialDynamoDBTable.jpg)

## Using
Once deployed and substituting your `<base URL>` the following CURL commands can be used to interact with the resulting API, whose results can be confirmed in the DynamoDB console

### POST

```bash
curl -X POST https:<base URL>/go-sls-crudl -d @data/post1.json
```
Which should result in the DynamoDB table looking like this:

![First Post DynamoDB Table](/img/firstPostDynamoDBTable.jpg)

Rinse/repeat for other data files to yeild:

![All Posts DynamoDB Table](/img/allPostsDynamoDBTable.jpg)

### GET Specific Item
Using the year and title (replacing spaces wiht '-' or '+'), you can now obtain an item as follows (prettified output):
```bash
curl https://<base URL>/go-sls-crudl/2013/Hunger-Games-Catching-Fire
{
  "year": 2013,
  "title": "Hunger Games Catching Fire",
  "info": {
    "plot": "Katniss Everdeen and Peeta Mellark become targets of the Capitol after their victory in the 74th Hunger Games sparks a rebellion in the Districts of Panem.",
    "rating": 7.6
  }
}
```

### GET a List of Items
You can list items by year as follows (prettified output):
```bash
curl https://<base URL>/go-sls-crudl/2013
[
  {
    "year": 2013,
    "title": "Hunger Games Catching Fire",
    "info": {
      "plot": "",
      "rating": 0
    }
  },
  {
    "year": 2013,
    "title": "Turn It Down Or Else",
    "info": {
      "plot": "",
      "rating": 0
    }
  }
]
```

### DELETE Specific Item
Using the same year and title specifiers, you can delete as follows:
```bash
curl -X DELETE https://<base URL>/go-sls-crudl/2013/Hunger-Games-Catching-Fire
```
Which should result in the DynamoDB table looking like this:

![First Delete DynamoDB Table](/img/firstDeleteDynamoDBTable.jpg)

### UPDATE Specific Item
You can update as follows:
```bash
curl -X PUT https:<base URL>/go-sls-crudl -d @data/put3.json
```
Which should result in the DynamoDB table looking like this:

![First Update DynamoDB Table](/img/firstUpdateDynamoDBTable.jpg)
