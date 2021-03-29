# bbolt API

This is a simple REST API for interacting with bbolt. The server looks for an environment variable named `DATABASE_PATH`. 

The server only opens one database. The server is lightweight, so if you need more than one database you can just run
another server.

## Buckets

`/`