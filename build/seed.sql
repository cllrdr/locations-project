-- Seed data для locations-project
ALTER DATABASE "locations-project" SET timezone TO 'Europe/Moscow';
SELECT pg_reload_conf();

-- Локации
INSERT INTO locations (id, name, description, image_path, video_path, players, is_deleted) VALUES
(1, 'Тирсфальские леса', 'Заброшенные леса Лордерона с богатыми залежами золота и лесными угодьями.', 'http://localhost:9000/locations/map1.png', NULL, '3-6', FALSE),
(2, 'Пустоши Дурхота', 'Сухие красные степи, где золото встречается часто, а дерево можно найти только в редких оазисах.', 'http://localhost:9000/locations/map2.png', NULL, '2-4', FALSE),
(3, 'Ледяная Корона', 'Суровая мерзлота с незначительными запасами древесины, но богатыми золотыми жилами глубоко во льдах.', 'http://localhost:9000/locations/map3.png', NULL, '4-8', FALSE),
(4, 'Болота Печали', 'Топкое негостеприимное место, бедное на любые ресурсы — только выживание и контроль над ограниченными источниками.', 'http://localhost:9000/locations/map4.png', NULL, '1-3', FALSE),
(5, 'Пылающие степи', 'Выжженная демоническая земля, где нет деревьев, но золото течет рекой (в прямом смысле — жидкое золото в лаве).', 'http://localhost:9000/locations/map5.png', NULL, '2-4', FALSE);

-- Пользователи
INSERT INTO users (id, email, name, password, is_moderator) VALUES
(1, 'alpha@example.com', 'Player_Alpha', 'password123', FALSE),
(2, 'beta@example.com', 'Player_Beta', 'password123', FALSE),
(3, 'gamma@example.com', 'Player_Gamma', 'password123', FALSE);

-- Заявки игроков (черновик для пользователя 1)
INSERT INTO players_location_requests (id, nickname, status, created_at, creator_id, moderator_id) VALUES
(1, 'Player_Alpha', 'черновик', NOW(), 1, NULL);

-- Выбранные локации (добавим одну локацию в черновик)
INSERT INTO players_chosen_locations (request_id, location_id, priority) VALUES
(1, 1, 1);
