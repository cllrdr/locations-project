-- Seed data для locations-project
ALTER DATABASE "locations-project" SET timezone TO 'Europe/Moscow';
SELECT pg_reload_conf();

-- Локации
INSERT INTO locations (name, description, short_description, image_path, video_path, players, is_deleted) VALUES
('Тирсфальские леса', 'Заброшенные леса Лордерона с богатыми залежами золота и лесными угодьями.', 'Ancient forest with gold deposits', 'http://localhost:9000/locations/map1.png', 'http://localhost:9000/locations/map1.mp4', '3-6', FALSE),
('Пустоши Дурхота', 'Сухие красные степи, где золото встречается часто, а дерево можно найти только в редких оазисах.', 'Red steppes and oases', 'http://localhost:9000/locations/map2.png', 'http://localhost:9000/locations/map2.mp4', '2-4', FALSE),
('Ледяная Корона', 'Суровая мерзлота с незначительными запасами древесины, но богатыми золотыми жилами глубоко во льдах.', 'Icy lands with golden veins', 'http://localhost:9000/locations/map3.png', 'http://localhost:9000/locations/map3.mp4', '4-8', FALSE),
('Болота Печали', 'Топкое негостеприимное место, бедное на любые ресурсы — только выживание и контроль над ограниченными источниками.', 'Dangerous swamps', 'http://localhost:9000/locations/map4.png', 'http://localhost:9000/locations/map4.mp4', '1-3', FALSE),
('Пылающие степи', 'Выжженная демоническая земля, где нет деревьев, но золото течет рекой (в прямом смысле — жидкое золото в лаве).', 'Demonic lava fields', 'http://localhost:9000/locations/map5.png', 'http://localhost:9000/locations/map5.mp4', '2-4', FALSE);

-- Пользователи
INSERT INTO users (email, name, password, is_moderator) VALUES
('alpha@example.com', 'Player_Alpha', 'password123', FALSE),
('beta@example.com', 'Player_Beta', 'password123', FALSE),
('gamma@example.com', 'Player_Gamma', 'password123', FALSE);

-- Заявки игроков (черновик для пользователя 1)
INSERT INTO players_location_games (nickname, status, created_at, creator_id, moderator_id) VALUES
('Player_Alpha', 'черновик', NOW(), 1, NULL);

-- Выбранные локации (добавим одну локацию в черновик)
INSERT INTO players_chosen_locations (request_id, location_id, priority) VALUES
(1, 1, 1);
