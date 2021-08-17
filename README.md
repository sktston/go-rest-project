# go-rest-project

### Run
```
$ go run main.go
```

### Generate swagger files
```
$ swag init
```

It will create files in `/docs` by reading swagger-related comments. (Refer `/handler/book_handler.go`)
The generated swagger json is loaded when the server boots.
After the server boots, visit http://localhost:8080/swagger/index.html.