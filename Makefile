build: tidy format
	go build

serve: tidy format
	go run .

live:
	rsync -rtvzP ./htdocs/ sam@samanthony.xyz:/var/www/htdocs/samanthony.xyz/

format:
	gofmt -s -w .

tidy:
	go mod tidy
