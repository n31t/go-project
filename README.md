# Anime API

This is a simple RESTful API built in Go that allows you to manage an anime database. It uses the Gorilla Mux router and JSON for request and response bodies. The idea is that we have several entities, such as User, Anime, WatchedAnime, AnimeForTierList.

## Endpoints

- `POST /animes`: Create a new anime. The request body should include title, episodes, studio, description, releaseYear, genre, and rating.

- `GET /animes`: Get a list of all animes.

- `GET /animes/{id}`: Get the details of a specific anime by its ID.

- `PUT /animes/{id}`: Update the details of a specific anime by its ID. The request body can include title, episodes, studio, description, releaseYear, genre, and rating.

- `DELETE /animes/{id}`: Delete a specific anime by its ID.

## Error Handling

The API responds with appropriate HTTP status codes and error messages in the case of an error. For example, if you try to get an anime that doesn't exist, you'll receive a 404 status code and an "Anime not found" message.

## Usage

To start the server, run `go run main.go`. The server will start on port 8081. You can then use a tool like curl or Postman to send requests to the API.

## Future Improvements

- Add authentication to protect sensitive endpoints.
- Add more detailed error messages for debugging purposes.
- Add endpoints for User, WatchedAnime, AnimeForTierList.
