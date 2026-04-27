-- Локации
INSERT INTO locations (name, description, image_path, video_path, players) VALUES
('Тирсфальские леса', 'Заброшенные леса Лордерона с богатыми залежами золота и лесными угодьями.', 'http://localhost:9000/locations/map1.jpg', 'http://localhost:9000/locations/map1.mp4', '3-6'),
('Пустоши Дурхота', 'Сухие красные степи, где золото встречается часто, а дерево можно найти только в редких оазисах.', 'http://localhost:9000/locations/map2.jpg', 'http://localhost:9000/locations/map2.mp4', '2-4'),
('Ледяная Корона', 'Суровая мерзлота с незначительными запасами древесины, но богатыми золотыми жилами глубоко во льдах.', 'http://localhost:9000/locations/map3.jpg', 'http://localhost:9000/locations/map3.mp4', '4-8'),
('Болота Печали', 'Топкое негостеприимное место, бедное на любые ресурсы — только выживание и контроль над ограниченными источниками.', 'http://localhost:9000/locations/map4.jpg', 'http://localhost:9000/locations/map4.mp4', '1-3'),
('Пылающие степи', 'Выжженная демоническая земля, где нет деревьев, но золото течет рекой (в прямом смысле — жидкое золото в лаве).', 'http://localhost:9000/locations/map5.jpg', 'http://localhost:9000/locations/map5.mp4', '2-4');

-- Пользователи
INSERT INTO users (email, name, password, is_moderator) VALUES
('alpha@example.com', 'Player_Alpha', 'password123', false),
('beta@example.com', 'Player_Beta', 'password123', false),
('gamma@example.com', 'Player_Gamma', 'password123', false);

