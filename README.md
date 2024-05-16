# Anime API

This is a simple RESTful API built in Go that allows you to manage an anime database. It uses the Gorilla Mux router and JSON for request and response bodies. The idea is that we have several entities, such as User, Anime, WatchedAnime, AnimeForTierList.

## Endpoints

### Animes
- `POST /api/v1/animes`: Create a new anime. The request body should include title, episodes, studio, description, releaseYear, genre, and rating.
- `GET /api/v1/animes`: Get a list of all animes.
- `GET /api/v1/animes/{id}`: Get the details of a specific anime by its ID.
- `PUT /api/v1/animes/{id}`: Update the details of a specific anime by its ID. The request body can include title, episodes, studio, description, releaseYear, genre, and rating.
- `DELETE /api/v1/animes/{id}`: Delete a specific anime by its ID.

### Watched Animes
- `POST /api/v1/watched-animes`: Create a new watched anime entry.
- `GET /api/v1/watched-animes/{id}`: Get the details of a specific watched anime by its ID.
- `GET /api/v1/watched-animes`: Get a list of all watched animes.
- `GET /api/v1/watched-animes/tier/{tier}`: Get a list of all watched animes by tier.
- `PUT /api/v1/watched-animes/{id}`: Update the details of a specific watched anime by its ID.
- `DELETE /api/v1/watched-animes/{id}`: Delete a specific watched anime by its ID.

### Users
- `POST /api/v1/users`: Register a new user.
- `PUT /api/v1/users/activated`: Activate a user.

### Tokens
- `POST /api/v1/tokens/authentication`: Create an authentication token.

### Healthcheck
- `GET /api/v1/healthcheck`: Check the health of the API.

## Error Handling

The API responds with appropriate HTTP status codes and error messages in the case of an error. For example, if you try to get an anime that doesn't exist, you'll receive a 404 status code and an "Anime not found" message.

## Usage

To start the server, run `go run *`. The server will start on port 8081. You can then use a tool like curl or Postman to send requests to the API.

## Deployment

This project is currently deployed on [https://neit-project-ple6u.ondigitalocean.app/api/v1/healthcheck](DigitalOcean). You can access the API using the provided endpoint.

## Environment Variables

For the API to work properly, make sure to provide the required environment variables in a .env file. Here's an example of the required variables:

- DSN=postgres://postgres:password@postgres:5432/adilovamir?sslmode=disable
- user=postgres
- dbname=adilovamir
- password=password
- host=db
