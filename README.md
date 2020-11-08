# HearMeMail

`HearMeMail` is a basic server that allows the user to send emails via Http calls.

## Setup

### Running the server using docker

1. Create network bridge for communication between containers:

        docker network create hearme-bridge
2. Start Postgres (if you have not pulled the postgres docker image follow the steps `Running Postgres in Docker``):

        docker run --name hearme-postgres -e POSTGRES_PASSWORD=<your postgres password> --network hearme-bridge -p 5432:5432 -d postgres

3. Build docker image

        docker build -t hearme-mail:latest .

4. Start docker container

        docker run --network hearme-bridge --name hearme-mail -p 8080:8080 hearme-mail

### Running server without docker

1. Clone Repository:

        git clone https://github.com/HugoGreyvenstein/hearMeMail.git
2. Change directory to project root:

        cd hearMeMail

3. Edit the configuration file (default location is the project root: `config.yml`):

        email:
          api-key: "<your api key>"
          name: "<your name>"
          from: "<your sendgrid email address>"

4. Start server:
    
        go run main.go
 The server will be started on port `8080`
        
The default config file is in the project root and is named `config.yml`. 
    A custom config file can be specified by supplying a command line argument:
    
        go run main.go <config-file-location>

### Running Postgres in Docker

1. Pull image from DockerHub (https://hub.docker.com/_/postgres):

        docker pull postgres
2. Run image as container:

        docker run --name hearme-postgres -e POSTGRES_PASSWORD=<your preferred password> -p 5432:5432 -d postgres
    - `--name`: human recognisable name for your container
    - `-e`: sets environment variables
    - `POSTGRES_PASSWORD`: root postgres password
    - `-p`: connects the port running inside of the container with your computer's port i.e. `<container port>:<your computer's port>`
    - `-d`: specifies that the command will return while the container runs independently
3. Check your postgres connection:
    Postgres can now be reached via port `5432`. Use your DB client of choice to test the connection.

## Troubleshooting

Some dependencies my need to be downloaded. 
    The following command will download the required dependencies. 
    Make sure you are in the project root directory:
    
        go get -v ./...
        
If the following error occurrs there is already a `HearMeMail` server running.
    Find the running process and kill it:
    
        Only one usage of each socket address (protocol/network address/port) is normally permitted.

## Future improvements

- How the ORM (Gorm) is used can be improved upon. At the moment transactions and concurrent calls may be handled incorrectly.
- Refactor to use a logger such as Logrus
- Make config file configurable
- Refactor docker image building so that config can be changed without building a new docker image
- Refactor builder constructors into fluent methods
- HTTP level error handling can be generalised. At the moment a lot of the error handling is duplicated between handlers
- Database should use a connection pool
- Docker layers should cache dependencies to avoid downloading when image is built
- Add indices on database tables