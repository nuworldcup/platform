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

Now when the server application is started, it will perform the necessary database migrations found in the `migrate/migrations` directory, using the `migrate.go` script.

To see that the first migration up has worked you can check the db to see that the changes are there:

```
psql nuwc -c "\d player"
```

- if you are in need of a Postgres client, [Postico](https://eggerapps.at/postico/) is good!

For basic `psql` commands visit this [cheatsheet](https://jazstudios.blogspot.com/2010/06/postgresql-login-commands.html)

#### Add Tournaments

In order for someone to use most of the endpoint, it assumes that there is a tournament made. In the current documentation

TODO:: Endpoint that creates tournaments so admins can create this from a web or mobile app

```
INSERT INTO tournament(tournament_id, tournament_name, tournament_timestamp, registration_timestamp) VALUES('coed2020', '2020 Northwestern World Cup', 'May 22 09:00:00 2020 CST', 'February 9 13:00:00 2020 CST');

INSERT INTO tournament(tournament_id, tournament_name, tournament_timestamp, registration_timestamp) VALUES('womens2020', '2020 Women''s NU World Cup', 'May 22 09:00:00 2020 CST', 'March 9 13:00:00 2020 CST');
```

#### Google Sheets setup

For the app to be able update a spreadsheet you need to add your spreadsheet ID in `server.go`. Do not push `token.go` to github as it contains a token only acquired after authorization from a google account. If you want to use google sheets functionality locally, you'll need to manually authenticate on a browser the first time you use one of the available functions in `sheets.go`. Follow the instructions on the command line and in the browser and a `token.json` file will be added to `gtools/` and all api calls from that point forward will use the `token.json` file to authenticate.

## Run the server

Navigate to the root directory and start the server:

```
go run server.go
```

## Making structural db changes

If you're adding a table, index, adding a column or whatever it is, add another migration file. Make sure you put the sequential number in front of it, and add both a migrate up and migrate down. If you need help understanding migrations, it might be a good idea to look at the [golang-migrate docs](https://github.com/golang-migrate/migrate).

You can use the [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) to help you. Basically if you make a mistake, you can [force a version](https://github.com/golang-migrate/migrate/issues/282#issuecomment-530732246). Or if you just want to see what a prior version of the db is like you can migrate up or down accordingly. To install the golang-migrate CLI on mac:

```
brew install golang-migrate
```

If you want to work with a database version that doesn't include the latest version of migrations, comment out `migrate.Migrate()` in `main()` of `server.go()`, then use the go-lang CLI to migrate to whatever verison you want to use and restart the application. Optionally, you could alter `migrate.go`, but don't push those changes.

## Error Handling

If returning repsonses to a client, you can use `MalformedRequest` in the `/lib` package to create an error to return. This way if an error occurs, we check if it's of type `MalformedRequest` (we created the error and it's not a system issue) and then we return the http status code and message that corresponds to the error that the user's mistake. If it is not of that type, we can return an Internal System Error because we don't know what kind of issue it is.

## System Diagram

    ___________________________        __________________________________________
    | client (react front-end) | <--> | backend (platform) server <--> db server |
    ---------------------------        ------------------------------------------

    _____________________________________        ______________________________________       ___________
    | client (react front-end) on heroku | <--> | backend (platform) server on heroku | <--> | db server |
    -------------------------------------        --------------------------------------       -----------
