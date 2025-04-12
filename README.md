# The sequel

Despite the name, the sequel is an in memory key value database built in Go

## Features

1. tcp server listener to recieve instructions
2. database persistance
3. thread-safe mutations and queries

## Next steps

1. Implement a more efficient thread-safe pattern

## Building & Running the database

To build the engine run the following command

```sh
go build ./cmd/server/main.go
```

Then the database can be ran by calling the executable that was just built
