# webservice
Go Minimal Webservice Example with Test cases and Coverage Reports


# Introduction

This is a minimal project replicating Rest APIs with Test cases and coverage reports. 
Currently using Regex for routing and using in-memory storage

Here is a table showing sample requests

|   	**Url**|   	**Method**|   	**Request Body**|
|:-:	|:-:	|:-:	|
|   	/users/|   	GET|   	-|
|   	/users/|   	POST|   	```{ "FirstName": "John", "LastName": "Zick" }```|
|   	/users/|   	PUT|   	```{ "Id": 1,"FirstName": "John", "LastName": "Wick" }```|
|   	/users/1|   	GET|   	-|
|   	/users/1|   	DELETE|   	-|

# Test Coverage

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


