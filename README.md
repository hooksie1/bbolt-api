# bbolt API

This is a simple REST API for interacting with bbolt. The server looks for an environment variable named `DATABASE_PATH`
 and an environment variable named `SERVER_PORT`.

The server only opens one database. The server is lightweight, so if you need more than one database you can just run
another server.

## API Documentation

Visit the [Docs](https://docs.bbolt-api.hooks.technology) site

The application now has an API explorer built in using the Swagger UI. It's located at `/v1/docs/`.

## Running With Docker

`docker run -d --name api -v bolt-volume:/database -p 8080:8080 hooksie1/bbolt-api:v0.0.4`

The container automatically stores the database named `bolt.db` in `/database`. Mount a volume 
to `/database` to store locally. 
