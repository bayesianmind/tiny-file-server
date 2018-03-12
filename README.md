# Tiny File Server
This is a simple tiny file server in Golang that implements basic user
management and user specific file fetching capabilities.

It is currently using in-memory user management and on-disk file storage
 but could easily be extended to work with blob storage and a durable datastore.

## Running
```
dep ensure
go run cmd/main.go
```

## Tests
```
go test ./...
```