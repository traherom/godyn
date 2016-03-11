= Building
Tested with Go 1.6. A simple `go get && go build` should work.

To make it extra clean, the integrated `build.sh` and `run.sh` scripts will
create Docker containers to run go dyn.

= Parameters
For ease of using it in Docker, godyn pulls its parameters from environment
variables: 

    * GODYN_SERVICE = Domain service you're using (`domains.google.com`)
    * GODYN_HOST = Host name you're updating (`myip.example.com`)
    * GODYN_USER = User name of the service you're using
    * GODYN_PASSWORD = Password for the service

