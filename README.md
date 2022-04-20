# GRPC tests
This repository is a test of some grpc capabilities in go, using docker and other containerization technologies for maximum scalability.

## Compiling the proto files
Fist we need to compile all files that will define our grpc protocol, all coded in the file functions/functions.proto with the following command:

`cd functions && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative functions.proto && cd ..`

## Compiling the project locally
After that, the go files for the client and server can be compiled and run (in separate terminals) locally for testing the code itself and ensure everythig runs correctly:

`go build -o grpc_server ./server/`
`go build -o grpc_client ./client/`

After running both executables (first the sever then the client) it is possible to see one window receiving strings and the other receiving echos of those strings 

## Running with Docker
The project can be converted to docker images with the commands:

`docker build -f server/Dockerfile -t grpc_test/server .`
`docker build -f client/Dockerfile -t grpc_test/client .`

Creating a proper local network for both containers to communicate is necessary:

`docker network create grpc`

Running both containers is 

`docker run -d --net grpc --name grpc_server -p 50501:50501 grpc_test/server:latest`
`docker run --net grpc grpc_test/client:latest`

It should be seen the same as in the local case, however, as the server is running in detached mode to access it's logs:

`docker logs grpc_server`