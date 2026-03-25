-- noinspection SqlDialectInspectionForFile

-- Seed data для locations-project (для продакшена с внешним IP)
ALTER DATABASE "locations-project" SET timezone TO 'Europe/Moscow';
SELECT pg_reload_conf();
-- Локации
INSERT INTO locations (id, name, description, image_path, video_path, players) VALUES
(1, 'Тирсфальские леса', 'Заброшенные леса Лордерона с богатыми залежами золота и лесными угодьями.', 'http://192.168.1.6:9000/locations-images/map1.jpg', 'http://192.168.1.6:9000/locations-videos/map1.mp4', '3-6'),
(2, 'Пустоши Дурхота', 'Сухие красные степи, где золото встречается часто, а дерево можно найти только в редких оазисах.', 'http://192.168.1.6:9000/locations-images/map2.jpg', 'http://192.168.1.6:9000/locations-videos/map2.mp4', '2-4'),
(3, 'Ледяная Корона', 'Суровая мерзлота с незначительными запасами древесины, но богатыми золотыми жилами глубоко во льдах.', 'http://192.168.1.6:9000/locations-images/map3.jpg', 'http://192.168.1.6:9000/locations-videos/map3.mp4', '4-8'),
(4, 'Болота Печали', 'Топкое негостеприимное место, бедное на любые ресурсы — только выживание и контроль над ограниченными источниками.', 'http://192.168.1.6:9000/locations-images/map4.jpg', 'http://192.168.1.6:9000/locations-videos/map4.mp4', '1-3'),
(5, 'Пылающие степи', 'Выжженная демоническая земля, где нет деревьев, но золото течет рекой (в прямом смысле — жидкое золото в лаве).', 'http://192.168.1.6:9000/locations-images/map5.jpg', 'http://192.168.1.6:9000/locations-videos/map5.mp4', '2-4');

-- Пользователи
INSERT INTO users (id, email, name, password, is_moderator) VALUES
-- (1, 'alpha@example.com', 'Player_Alpha', 'password123', false),
(2, 'beta@example.com', 'Player_Beta', 'password123', false),
(3, 'gamma@example.com', 'Player_Gamma', 'password123', false);

-- Заявки игроков
-- INSERT INTO players_location_requests (id, nickname, status, created_at, creator_id) VALUES
-- (1, 'Player_Alpha', 'сформирован', NOW(), 1),
-- (2, 'Player_Beta', 'сформирован', NOW(), 2),
-- (3, 'Player_Gamma', 'сформирован', NOW(), 3);
--
-- -- Выбранные локации
-- INSERT INTO players_chosen_locations (request_id, location_id, priority) VALUES
-- -- Заявка #1: Player_Alpha
-- -- (1, 2, 1),
-- -- (1, 3, 2),
-- -- (1, 4, 3),
-- -- Заявка #2: Player_Beta
-- (2, 2, 1),
-- (2, 4, 2),
-- -- Заявка #3: Player_Gamma
-- (3, 5, 1);
