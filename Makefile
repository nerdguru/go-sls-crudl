go_apps = bin/get bin/post bin/delete bin/put bin/list-by-year

bin/% : functions/%.go functions/moviedao.go
	env GOOS=linux go build -ldflags="-s -w" -o $@ $< functions/moviedao.go

build: $(go_apps) | vendor

vendor: Gopkg.toml
	dep ensure
