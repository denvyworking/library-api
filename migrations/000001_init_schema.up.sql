-- Создаем таблицы
CREATE TABLE IF NOT EXISTS authors (
    id SERIAL PRIMARY KEY,
    author VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS genres (
    id SERIAL PRIMARY KEY,
    genre VARCHAR(50) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS books (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    author_id INTEGER REFERENCES authors(id),
    genre_id INTEGER REFERENCES genres(id),
    price INTEGER NOT NULL CHECK (price >= 0)
);