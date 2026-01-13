-- Очищаем тестовые данные
DELETE FROM books WHERE id IN (1, 2, 3, 4, 5);
DELETE FROM authors WHERE id IN (1, 2, 3);
DELETE FROM genres WHERE id IN (1, 2, 3, 4, 5);