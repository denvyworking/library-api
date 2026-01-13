-- Вставляем тестовые данные
INSERT INTO authors (author) VALUES
    ('Лев Толстой'),
    ('Федор Достоевский'),
    ('Антон Чехов')
ON CONFLICT DO NOTHING;

INSERT INTO genres (genre) VALUES 
    ('Роман'),
    ('Драма'),
    ('Комедия'),
    ('Трагедия'),
    ('Фантастика')
ON CONFLICT DO NOTHING;

INSERT INTO books (name, author_id, genre_id, price) VALUES 
    ('Война и мир', 1, 1, 500),
    ('Анна Каренина', 1, 1, 450),
    ('Преступление и наказание', 2, 1, 400),
    ('Братья Карамазовы', 2, 1, 420),
    ('Вишневый сад', 3, 2, 300)
ON CONFLICT DO NOTHING;