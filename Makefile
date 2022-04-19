build_dev: tidy format
	go build -o devserver

serve_dev: build_dev
	./devserver --chroot ./ --root /htdocs/ --host localhost

build: tidy format
	GOOS=openbsd GOARCH=amd64 go build -o webserver

live:
	rsync -rtvzP ./htdocs/ sam@samanthony.xyz:/var/www/htdocs/samanthony.xyz/

format:
	gofmt -s -w .

tidy:
	go mod tidy
