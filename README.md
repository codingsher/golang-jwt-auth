This API includes complete implementation of Http Requests i.e Get, Post, Patch, Delete using Golang's net/http standard library and SQLite database. This API has been coded using production level API structure that is efficient for future plugin integration/changes, testing etc.
This is a simple RestAPI that uses SQLite to store students data i.e Name, Email, and Age.

To test the API, follow these steps:

1. Clone into your local repo.
2. go mod tidy
3. add your local.yaml file
4. go run cmd/server/main.go -config config/local.yaml
5. open postman to test API.
