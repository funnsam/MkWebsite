build-builder-windows:
	go env -w GOOS=windows GOARCH=amd64

	go mod tidy
	go build -o ./mkwebsite.exe ./src/

build-builder-unix:
	go env -w GOOS=linux GOARCH=amd64

	go mod tidy
	go build -o ./mkwebsite ./src/