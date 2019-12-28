CREATE TABLE events (
  time timestamp NOT NULL DEFAULT now(),
  chat_id numeric(15) NOT NULL,
  user_id numeric(15) NOT NULL,
  credits int NOT NULL,
  creator_id numeric(15) NOT NULL,
  telegram_update int NOT NULL
);

CREATE INDEX ON events(time);
CREATE INDEX ON events(chat_id, user_id);
