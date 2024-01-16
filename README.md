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

This project contains 1 initial migration, which is the `create_user_table` migration
