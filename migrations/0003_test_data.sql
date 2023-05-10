-- +goose Up
-- тестовые города
insert into bx_cities(name) values
    ('Краснодар'),
    ('Ульяновск');

-- тестовые отделы
insert into bx_departments(name, lead) values
   ('Backend Golang', null),
   ('Backend Java', null),
   ('Backend PHP', null),
   ('Backend PHP Bitrix', null),
   ('Backend Python', null),
   ('Frontend React.js', null),
   ('Frontend Vue.js', null),
   ('Mobile Android', null),
   ('Mobile React Native', null),
   ('Mobile Swift', null),
   ('Направление 1С', null),
   ('Направление AI/ML/DS', null),
   ('Направление DevOps', null);

-- тестовые пользователи
insert into bx_users(id, name, last_name, joined_at, leaved_at, department, city) values
    (1235, 'Игорь', 'Денисенко', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (0, 'Дмитрий', 'Смолов', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (1, 'Андрей', 'Смолов', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (2, 'Сергей', 'Косов', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (3, 'Амир', 'Курамов', '2023-03-10 19:40:18.000000', null, 'Backend Java', 'Краснодар'),
    (4, 'Инга', 'Долгая', '2023-03-10 19:40:18.000000', null, 'Направление AI/ML/DS', 'Краснодар'),
    (5, 'Андрей', 'Циров', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (6, 'Владимир', 'Оленин', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (7, 'Исав', 'Федов', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (8, 'Лев', 'Акимов', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (9, 'Егор', 'Зимин', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (5534, 'Дмитрий', 'Измайлов', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар'),
    (4321, 'Кирилл', 'Новоселецкий', '2023-03-10 19:40:18.000000', null, 'Backend Golang', 'Краснодар');

-- привязка техлидов к отделам
update bx_departments set lead = 4321 where name = 'Backend Golang';
update bx_departments set lead = 5534 where name = 'Backend Java';
update bx_departments set lead = 4321 where name = 'Backend PHP';
update bx_departments set lead = 4321 where name = 'Backend PHP Bitrix';
update bx_departments set lead = 4321 where name = 'Backend Python';
update bx_departments set lead = 5534 where name = 'Frontend React.js';
update bx_departments set lead = 4321 where name = 'Frontend Vue.js';
update bx_departments set lead = 4321 where name = 'Направление 1С';
update bx_departments set lead = 5534 where name = 'Направление AI/ML/DS';

-- тестовые чаты
insert into tg_chats (id, sort, type, title, description, joined_at, invite_link)
values  (-884525477, 100, 'all', 'ZeTest.ZeBrains', 'Общий чат, обязателен для всех', '2023-03-10 19:40:18.000000', 'https://t.me/+sRBs_sWN_Hs0Yjli'),
        (-1001839122979, 100, 'all', 'ZeTest.Info', 'Общий канал, обязателен для всех', '2023-03-10 19:54:00.000000', 'https://t.me/+DrlD9cVVQV81YTVi'),
        (-763702331, 100, 'department', 'ZeTest.Dev', 'Общий чат, обязателен для всех разработчиков', '2023-03-10 19:54:00.000000', 'https://t.me/+6WQTWJpCQiUzNTIy'),
        (-999564057, 200, 'department', 'ZeTest.Golang', 'Чат стека, обязателен для отдела "Backend Golang"', '2023-03-10 19:46:58.000000', 'https://t.me/+loGhdA864js1MTAy'),
        (-953306580, 100, 'city', 'ZeTest.Краснодар', 'Чат города, обязателен для Краснодара', '2023-03-10 19:55:52.000000', 'https://t.me/+z4v0KQOG_5hmODNi'),
        (-999353704, 100, 'recommended', 'ZeTest.Anime', 'Чат по интересам, необязателен', '2023-03-10 19:58:48.000000', 'https://t.me/+Q-Ui_NNtq90xY2Y6'),
        (-843962996, 100, 'other', 'ZeTest.NoAdmin', 'Чат без административных прав', '2023-03-10 19:58:48.000000', 'https://t.me/+L363VsHtj9JlM2Ji');

-- привязка чатов к городам
insert into bx_cities_chats (tg_chat, bx_city) values
    (-953306580, 'Краснодар');

-- привязка чатов к отделам
insert into bx_departments_chats (tg_chat, bx_department) values
    (-763702331, 'Backend Golang'),
    (-999564057, 'Backend Golang');

-- тестовые пользователи
insert into tg_users(id, nickname, bx_user) values
    (431414519, 'im_denisenko', 1235),
    (523019020, 'rust_env', 4321),
    (6072956963, 'test_dm', 5534);

-- привязка тестовых пользователей к чатам
insert into tg_chats_members (tg_user, tg_chat, joined_at) values
    (431414519, -884525477, '2023-03-10 19:55:52.000000'),
    (431414519, -763702331, '2023-03-10 19:55:52.000000'),
    (431414519, -999564057, '2023-03-10 19:55:52.000000'),
    (431414519, -953306580, '2023-03-10 19:55:52.000000'),
    (431414519, -999353704, '2023-03-10 19:55:52.000000'),
    (523019020, -884525477, '2023-03-10 19:55:52.000000'),
    (523019020, -763702331, '2023-03-10 19:55:52.000000'),
    (523019020, -999564057, '2023-03-10 19:55:52.000000'),
    (523019020, -953306580, '2023-03-10 19:55:52.000000'),
    (523019020, -999353704, '2023-03-10 19:55:52.000000'),
    (6072956963, -999564057, '2023-03-10 19:55:52.000000'),
    (6072956963, -953306580, '2023-03-10 19:55:52.000000'),
    (6072956963, -999353704, '2023-03-10 19:55:52.000000');
-- +goose Down