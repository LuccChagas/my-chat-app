CREATE TABLE "users" (
                         "id" uuid PRIMARY KEY,
                         "password" varchar NOT NULL,
                         "cpf" varchar NOT NULL,
                         "email" varchar NOT NULL,
                         "phone" varchar NOT NULL,
                         "name" varchar NOT NULL,
                         "first_name" varchar NOT NULL,
                         "last_name" varchar NOT NULL,
                         "nick_name" varchar NOT NULL,
                         "created_at" timestamptz NOT NULL DEFAULT (now()),
                         "updated_at" timestamp
);