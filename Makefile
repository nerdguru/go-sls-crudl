build:
	dep ensure
	env GOOS=linux go build -ldflags="-s -w" -o bin/get functions/get.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post functions/post.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete functions/delete.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/put functions/put.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/list-by-year functions/list-by-year.go
