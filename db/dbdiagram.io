Table notification.telegram_user {
  id int [pk]
  name varchar
  created_at timestamp
}

Table notification.notification {
  id int [pk, increment]
  user_id int [ref: > notification.telegram_user.id]
  value text
  schedule varchar
  created_at timestamp
  updated_at timestamp
  
  Indexes {
    (user_id) [name: 'user_id']
  }
}

SET client_min_messages TO WARNING;

CREATE SCHEMA IF NOT EXISTS "notification";

CREATE TABLE IF NOT EXISTS "notification"."telegram_user" (
  "id" int PRIMARY KEY,
  "name" varchar,
  "created_at" timestamp
);

CREATE TABLE IF NOT EXISTS "notification"."telegram_notification" (
  "id" INT GENERATED BY DEFAULT AS IDENTITY PRIMARY KEY,
  "user_id" int,
  "value" text,
  "schedule" varchar,
  "created_at" timestamp,
  "updated_at" timestamp
);

CREATE INDEX IF NOT EXISTS "user_id" ON "notification"."telegram_notification" ("user_id");
COMMENT ON COLUMN "notification"."telegram_notification"."user_id" 
	IS 'See notification.telegram_user.id';
COMMENT ON COLUMN "notification"."telegram_notification"."schedule" 
	IS 'Unix cron schedule format';
