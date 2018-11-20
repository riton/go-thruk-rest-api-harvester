all:
	CGO_ENABLED=0 GOOS=linux go build -o thruk-rest-api-harvester .