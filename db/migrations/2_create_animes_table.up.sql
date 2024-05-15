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

CREATE TABLE IF NOT EXISTS  watched_animes (
	id SERIAL PRIMARY KEY,
	user_id INT NOT NULL,
	anime_id INT NOT NULL,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (anime_id) REFERENCES animes(id),
	was_viewed TIMESTAMP(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
	tier TEXT DEFAULT 'none'
);