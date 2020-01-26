# platform

backend service for Northwestern World Cup

## Getting Started

This project needs to be in your go path. However, you can put this repo where ever you want on your local machine and create a symlink from your go path:

```
ln -s path/to/original your/go/path
```

You'll most likely want to create the link from `~/go/src` but it could be different if you don't use the defaults for go

#### Create the nuwc db

Install postgres on your local machine (best to do with brew) and let the server run on the default port `5432`. Create a user called `nuwcuser` with password `password` and create a db called `nuwc`. Grant all privelages for this db to the new user.

```
CREATE ROLE nuwcuser WITH LOGIN PASSWORD 'password';
CREATE DATABASE nuwc;
GRANT ALL PRIVILEGES ON DATABASE nuwc TO nuwcuser;
```

Now when the server application is started, it will perform the necessary database migrations found in the `migrate/migrations` directory.

To see that the first migration up has worked you can check the db to see that the changes are there:

```
psql nuwc -c "\d player"
```

- if you are in need of a Postgres client, [Postico](https://eggerapps.at/postico/) is good!

For basic `psql` commands visit this [cheatsheet](https://jazstudios.blogspot.com/2010/06/postgresql-login-commands.html)

## Run the server

Navigate to the root directory and start the server:

```
go run server.go
```
