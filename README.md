# webservice
Go Minimal Webservice Example with Http Request Logging (like tomcat access-log) & Request based logging levels, Test cases and Coverage Reports


## Introduction

This is a minimal project replicating Rest APIs with Http Request Logging (like tomcat access-log) & Request based logging levels, Test cases and coverage reports.

Currently using Mux & Regex for routing and using in-memory storage

Here is a table showing sample requests:

| **Url**| **Method**| **Request Body**|
|:-: |:-: |:-: |
| /users/| GET| -|
| /users/| POST| ```{ "FirstName": "John", "LastName": "Zick" }```|
| /users/1| PUT| ```{ "Id": 1,"FirstName": "John", "LastName": "Wick" }```|
| /users/1| GET| -|
| /users/1| DELETE| -|



Following are the CUrls for available operations:

## Post Request ##

Adds new user in application

```
curl --request POST \
--url http://localhost:3000/users \
--header 'Content-Type: application/json' \
--data '{
"FirstName": "John",
"LastName": "Zick",
"Debug": true
}'
```


## Get Request ##

Lists all users registered with application

```
curl --request GET \
--url http://localhost:3000/users \
--header 'Content-Type: application/json' \
--data '{
"Debug": true
}'
```


## Get By ID ##

Get single user by `id`

```
curl --request GET \
--url http://localhost:3000/users/6 \
--header 'Content-Type: application/json' \
--data '{
"Debug": true
}'
```


## Put Request ##

Updates existing user in application

```
curl --request PUT \
--url http://localhost:3000/users/1 \
--header 'Content-Type: application/json' \
--data '{
"Id": 1,
"FirstName": "John",
"LastName": "Wick",
"Debug": true
}'
```

## Delete Request ##

Deletes user by Id

```
curl --request DELETE \
--url http://localhost:3000/users/6 \
--header 'Content-Type: application/json' \
--data '{
"Debug": true
}'
```

**Note:** To change the logging level to default (info), just remove the request body from Http Request


# Working Explanation

## Access Logs ##

To generate access-logs (can be found in project dir as **access_logs.txt**), a middleware has been added by handler-function which logs on default (info) level which is held inside helper package as **logger.go** under `accessLogger *zap.SugaredLogger`


I have kept the middleware simple, but you can add application-level configurations to make it more efficient.


## Application Logs ##

The application logs (can be found in project dir as **logs.txt**), I have added `defaultLogger *zap.SugaredLogger` & `debugLogger *zap.SugaredLogger`. Both are held inside helper package, callers need to ask for logger by providing `http.Request` object where it makes the decision by extracting the following information from `Request Body`:

```
type Param struct {
Debug bool
}
```

So the usage from caller perspective will become:

```
var logger *zap.SugaredLogger = nil

func (controller UserController) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
logger = helpers.GetLoggerByRequest(request)
logger.Info("Inside User Controller Entrypoint")
logger.Debugw("User Controller Entrypoint", "URL.PATH", request.URL.Path, "Method", request.Method)
.....
```


Right now, I am holding the logger instance on package level, but I have plan to associate it with some mapping against current user to make it more consistent and support concurrency.


**.....**


# Brief Introduction on Test Coverage

Test coverage measures the percentage of your code that is executed during tests.

To generate a test coverage report, run the following command:


```
go test -cover
```

The output will show the test coverage percentage for each package and an overall summary.

You can save output to file by adding -o flag with desired filename like:

```
go test -coverprofile coverage.out
```

To view code coverage results in console.


```
go tool cover -func="coverage.out"
```

To view code coverage results in html page, use the following command:

```
go tool cover -html="coverage.out"
```

Here **coverage.out** is cover file generated from previous command.

Executing this command will open the default web browser with visual representation of Code Coverage



### Feel free to edit/expand/explore this repository

For feedback and queries, reach me on LinkedIn at [here](https://www.linkedin.com/in/usama28232/?original_referer=)