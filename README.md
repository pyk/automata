# automata

REST API server for curhat app project.

## Building a service

    $ docker build -t bayu/automata:v0 .

## Running a service
database service named `postgres` must running first.
    
    $ docker run -d --env-file="data.env" --name automata --link postgres:database bayu/automata:v0 automata