build: | vendor
	env GOOS=linux go build -ldflags="-s -w" -o bin/get functions/get.go functions/moviedao.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/post functions/post.go functions/moviedao.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete functions/delete.go functions/moviedao.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/put functions/put.go functions/moviedao.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/list-by-year functions/list-by-year.go functions/moviedao.go

vendor: Gopkg.toml
	dep ensure
