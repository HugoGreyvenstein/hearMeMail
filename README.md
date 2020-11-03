# HearMeMail

`HearMeMail` is a basic server that allows the user to send emails via Http calls.

## Setup

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



## Troubleshooting

Some dependencies my need to be downloaded. 
    The following command will download the required dependencies. 
    Make sure you are in the project root directory:
    
        go get -v ./...
        
If the following error occurrs there is already a `HearMeMail` server running.
    Find the running process and kill it:
    
        Only one usage of each socket address (protocol/network address/port) is normally permitted.
