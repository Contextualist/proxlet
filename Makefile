.PHONY: build clean aws gcf upload

build:
	env GOOS=linux go build -ldflags="-s -w" -o bin/p lambda.go gohttp.go

clean:
	rm -rf ./bin

upload:
	sls deploy --verbose

aws: clean build upload

gcf:
	sed -i '' "1s/main/proxlet/" gohttp.go
	gcloud functions deploy proxlet --runtime go111 --entry-point Handler --trigger-http
	sed -i '' "1s/proxlet/main/" gohttp.go
