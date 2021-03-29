# bbolt API

This is a simple REST API for interacting with bbolt. The server looks for an environment variable named `DATABASE_PATH`. 

The server only opens one database. The server is lightweight, so if you need more than one database you can just run
another server.

## Buckets

|Endpoint|Method|Action|
|--------|------|------|
|`/v1/buckets/{bucket}`| GET | Return if Bucket exists|
|`/v1/buckets/{bucket}`| POST | Create bucket|
|`/v1/buckets/{bucket}`| DELETE | Delete bucket |

## Keys
|Endpoint|Method|Action|
|--------|------|------|
|`/v1/buckets/{bucket}/keys` | GET | List all keys in bucket |
|`/v1/buckets/{bucket}/keys/{key}`| GET | Get key information |
|`/v1/buckets/{bucket}/keys/{key}` | POST | Create key |
|`/v1/buckets/{bucket}/keys/{key}` | DELETE | Delete key|


## Creating a  KV
POST to the endpoint `/v1/buckets/{bucket}/keys/{key}` with the data in a payload.

Ex:

`curl http://localhost:8080/v1/mybucket/keys/mykey -d '{"data": "myvalue"}'`