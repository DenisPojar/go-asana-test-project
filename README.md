# Go Asana Test Project

## Description
This project is an extractor that fetches data from the Asana API, specifically information related to users and projects. The data is periodically saved as JSON files.

---

## Prerequisites
Before running this project, you must have the following installed:

- [Go](https://golang.org/dl/) (version 1.20+ recommended)
- Access to an Asana API token

Verify Go is installed:

```bash
go version
```

## Set up your configuration
! See the **.env.example** file. You need to copy the values from it to a separate new **.env** file, and assign real values.

## Run the project

```
git clone https://github.com/your-username/go-asana-extractor.git
cd go-asana-extractor
go mod tidy
go run ./cmd
```

## Run the tests
```
cd api/v1
go test -v
```
