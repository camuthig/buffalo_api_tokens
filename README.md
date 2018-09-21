# Buffalo API Tokens

This project is to better acquaint myself with Go, [Buffalo](http://gobuffalo.io) and new features such as modules. 
It also serves as a place to put together a sample of a simple way to implement API tokens in a fairly secure manner.

## Token Design

The design of the API tokens in this project are meant to be highly convenient for end users, but allow for a
reasonable amount of security.

Upon signing in, a user will be dispatched an access token and a refresh token. The access token is a standard JWT,
containing an expiration time and the ID of the user as the subject claim. The expiration time on the token is 1 hour
but can be tweaked based on the needs of a project. The refresh token is a cryptographical secure string, concatenated
with a representation of the user's ID to ensure uniqueness across users.

Upon expiration of an access token, a user can send their refresh token to the server to receive both a new access
token and refresh token. Upon use of the refresh token, it is removed from the system and is no longer valid.

When logging out, a user should be capable of sending their refresh token to the server for deletion. A user should
also be capable of logging out of all devices at once, essentially removing all active refresh tokens.

## Setup

A Docker Compose file is provided that will boot up a Postgres database for testing. To start this, just run

```
docker-compose up -d
```

On first starting the server, the migrations will need to be run with

```
buffalo db migrate
```

From there, the Buffalo server can be run in development with

```
buffalo dev
```
