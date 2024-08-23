hello:
	echo "hello"
test:
	echo "testing application"
run:
	echo "running application"
	go run cmd/crack/main.go test.csv 5
build:
	echo "bulding application"
tidy:
	go mod tidy