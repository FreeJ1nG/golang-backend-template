This is a [Golang](https://go.dev/) template with built-in authentication and examples

## Getting started

First, copy the content in `.env.sample` to a new `.env` file

```
ENV=local
SERVER_PORT=9090

APP_NAME=[Enter your app name here]

DB_USER=postgres
DB_USER_PASSWORD=[Enter your password for your db user here]

DB_HOST=localhost
DB_NAME=some_db_name
DB_PORT=[The port which the db will expose to]
DB_PASSWORD=[Password for db]

JWT_SECRET_KEY=[Enter a secret string for producing and encrypting the JWT]
JWT_EXPIRY_IN_HOURS=[Amount in hours specifying the expiry time of the JWT]
```

Then, compose the docker that will run the database server on port specified on `.env`

```
make up
```

or

```
docker compose up -d
```

Lastly, to run the development server

```
make dev
```

or

```
air
```

## Interfaces

Golang does not allow import cycles, to counter that we define interfaces for each **Handler, Service, Repository and Util** in`app/interfaces`. See `app/interfaces/auth.go` for some example, any other reference to another module's instance will use this `interfaces.SomeInstance` interface type

For example, if a module called `problem` uses any sort of authentication, say `AuthService`, we will define the `authService` typing as `interfaces.AuthService` so that the problem module does not have any dependency to the auth module

## DTO (Data Transfer Object)

Any definition of request types and response types will be defined on `app/dto/auth.go`

## Utility

- `config.go` extracts the `.env` file and initialized a `Config` object with the datas extracted from the environment file
- `logger.go` this middleware logs any request that goes in the server
- `parser.go` this utility file handles any sort of request parsing and response encoding
- `route-protector.go` creates a wrapper for any route that will be protected, ensuring that any handler function that is wrapped has a user context passed along with the request object

## Migrations

A migration is a series of changes to a database (be it of the table, of the schema, or anything related to the database)

A migration consists of an up migration and a down migration, alongside a sequence id of the migration

To create a migration, run

```
make migration [name_of_migration]
```

This will then create an up migration file and down migration file, like so:

```
db/migrations/000001_create_user_table.down.sql
db/migrations/000001_create_user_table.up.sql
```

The up migration file specifies how the database should handle an update/change in the database, and the down migration file specifies how to undo said changes.

This project contains 1 initial migration, which is the `create_user_table` migration, this migration creates the initial user table with `id`, `username`, `firstName`, `lastName` and `password_hash` fields.

To apply unapplied migration(s), do

```
make migrate
```

To undo migration(s), do

```
make migrate-down
```

## Handler, Service, Repository, Util files

- Handler files are essentially controllers, these handlers will handle incoming requests to certain urls (each request object can be accessed through `r *http.Request`)
- Service files handle the bussiness logic part of the application
- Repository files handle any sort of database mutation, typically used by service files
- Util files handle any sort of additional logic that makes use of repository or service files

## Route Protector Wrapper

To protect certain routes to only signed in users (requests in which the Authorization header is filled with a Bearer token), use the `routeProtector.Wrapper` as can be seen in the `app/injections.go` file.
