# property-service

## How to Run

make sure your config.yaml is valid and pointing to a running mongoDB
```shell
docker run --name mongodb -p 27017:27017  mongodb/mongodb-community-server:latest
```

to run the server:
```shell
go run main.go
```

if your config file is in a different path then `./config.yaml` you can use `-c`/`--config` to provide the path to the config file

```shell
go run main.go --config <path to config file>
```

## Some things I did do

I designed this service as a REST backend with MongoDB, splitting the code into 3 levels:

1. Transport - the REST controller, handling contracts, input validation and some of the pagination
2. Logic - the property handler itself, handling the logic and calls to the database
3. Database - the state itself, handling filters

My thinking was that switching to GRPC or to an SQL database could be done without major changes to the business logic itself.

I also threw in a quick and easy CI to run the tests 

I chose to write tests only for the logic layer, and only on the exported functions. I tried to generate input where applicable, and used coverage to find flows that were not tested.


## Some things I did not do

For the sake of time, I purposefully neglected some parts of this service:
* any config for the REST server itself
* better error context (I left the error messages the IDE suggested)
* a proper contract, either swagger or (preferably) protobuf
* I didn't take care of any overflow issues, for convenience
* **I did not take care of transactions - some of the functions do more than one db operation and they must be done in a transaction**

