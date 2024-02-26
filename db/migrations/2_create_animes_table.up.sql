CREATE TABLE IF NOT EXISTS  animes (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    episodes INT NOT NULL,
    studio TEXT NOT NULL,
    description TEXT NOT NULL,
    releaseYear INT NOT NULL,
    genre TEXT NOT NULL,
    rating FLOAT NOT NULL
);
