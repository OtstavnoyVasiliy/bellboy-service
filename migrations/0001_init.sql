-- +goose Up
-- Types of chats
CREATE TYPE "chat_types" AS ENUM (
  'all',
  'city',
  'department',
  'recommended',
  'other'
);

CREATE TABLE "tg_chats" (
    "id" bigint PRIMARY KEY,
    "sort" bigint NOT NULL DEFAULT 500,
    "type" chat_types NOT NULL,
    "title" varchar NOT NULL,
    "joined_at" timestamp NOT NULL,
    "description" text NOT NULL,
    "invite_link" varchar NOT NULL
);

CREATE TABLE "tg_users" (
    "id" bigint PRIMARY KEY,
    "bx_user" bigint,
    "nickname" varchar NOT NULL
);

CREATE TABLE "tg_chats_members" (
    "tg_chat" bigint,
    "tg_user" bigint,
    "joined_at" timestamp NOT NULL,
    PRIMARY KEY ("tg_chat", "tg_user")
);

CREATE TABLE "bx_cities" (
    "name" varchar PRIMARY KEY
);

CREATE TABLE "bx_cities_chats" (
    "tg_chat" bigint,
    "bx_city" varchar,
    PRIMARY KEY ("tg_chat", "bx_city")
);

CREATE TABLE "bx_departments" (
    "name" varchar PRIMARY KEY,
    "lead" bigint
);

CREATE TABLE "bx_departments_chats" (
    "tg_chat" bigint,
    "bx_department" varchar,
    PRIMARY KEY ("tg_chat", "bx_department")
);

CREATE TABLE "bx_users" (
    "id" bigint PRIMARY KEY,
    "name" varchar NOT NULL,
    "last_name" varchar NOT NULL,
    "joined_at" timestamp NOT NULL,
    "leaved_at" timestamp,
    "city" varchar,
    "department" varchar
);

COMMENT ON COLUMN "tg_chats"."id" IS 'ID чата в Telegram';

COMMENT ON COLUMN "tg_chats"."sort" IS 'Сортировка';

COMMENT ON COLUMN "tg_chats"."type" IS 'Тип чата';

COMMENT ON COLUMN "tg_chats"."title" IS 'Название';

COMMENT ON COLUMN "tg_chats"."joined_at" IS 'Дата добавления бота в чат';

COMMENT ON COLUMN "tg_chats"."description" IS 'Описание';

COMMENT ON COLUMN "tg_chats"."invite_link" IS 'Инвайт ссылка';

COMMENT ON COLUMN "tg_users"."id" IS 'ID пользователя Telegram';

COMMENT ON COLUMN "tg_users"."bx_user" IS 'ID пользователя Битрикс24';

COMMENT ON COLUMN "tg_users"."nickname" IS 'Никнейм пользователя Telegram';

COMMENT ON COLUMN "tg_chats_members"."tg_chat" IS 'ID чата Telegram';

COMMENT ON COLUMN "tg_chats_members"."tg_user" IS 'ID пользователя Telegram';

COMMENT ON COLUMN "tg_chats_members"."joined_at" IS 'Дата вступления пользователя в чат';

COMMENT ON COLUMN "bx_cities"."name" IS 'Название';

COMMENT ON COLUMN "bx_cities_chats"."tg_chat" IS 'ID чата Telegram';

COMMENT ON COLUMN "bx_cities_chats"."bx_city" IS 'Город';

COMMENT ON COLUMN "bx_departments"."name" IS 'Название';

COMMENT ON COLUMN "bx_departments"."lead" IS 'Техлид';

COMMENT ON COLUMN "bx_users"."id" IS 'ID в Битрикс24';

COMMENT ON COLUMN "bx_users"."name" IS 'Имя';

COMMENT ON COLUMN "bx_users"."last_name" IS 'Фамилия';

COMMENT ON COLUMN "bx_users"."joined_at" IS 'Дата приёма на работу';

COMMENT ON COLUMN "bx_users"."leaved_at" IS 'Дата увольнения';

COMMENT ON COLUMN "bx_users"."city" IS 'Город';

COMMENT ON COLUMN "bx_users"."department" IS 'Отдел';

ALTER TABLE "tg_users" ADD FOREIGN KEY ("bx_user") REFERENCES "bx_users" ("id");

ALTER TABLE "tg_chats_members" ADD FOREIGN KEY ("tg_chat") REFERENCES "tg_chats" ("id");

ALTER TABLE "tg_chats_members" ADD FOREIGN KEY ("tg_user") REFERENCES "tg_users" ("id");

ALTER TABLE "bx_cities_chats" ADD FOREIGN KEY ("tg_chat") REFERENCES "tg_chats" ("id");

ALTER TABLE "bx_cities_chats" ADD FOREIGN KEY ("bx_city") REFERENCES "bx_cities" ("name");

ALTER TABLE "bx_departments" ADD FOREIGN KEY ("lead") REFERENCES "bx_users" ("id");

ALTER TABLE "bx_departments_chats" ADD FOREIGN KEY ("tg_chat") REFERENCES "tg_chats" ("id");

ALTER TABLE "bx_departments_chats" ADD FOREIGN KEY ("bx_department") REFERENCES "bx_departments" ("name");

ALTER TABLE "bx_users" ADD FOREIGN KEY ("city") REFERENCES "bx_cities" ("name");

ALTER TABLE "bx_users" ADD FOREIGN KEY ("department") REFERENCES "bx_departments" ("name");

-- +goose Down
DROP TABLE tg_chats;
DROP TABLE tg_users;
DROP TABLE tg_chats_members;
DROP TABLE bx_cities;
DROP TABLE bx_cities_chats;
DROP TABLE bx_departments;
DROP TABLE bx_departments_chats;
DROP TABLE bx_users;