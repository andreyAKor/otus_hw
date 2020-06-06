-- +goose Up
-- +goose StatementBegin
CREATE TABLE events (
	id serial NOT NULL,
	title varchar(255) NOT NULL,
	"date" timestamptz NOT NULL,
	duration int8 NOT NULL,
	descr text NULL,
	user_id int4 NOT NULL,
	duration_start int8 NULL,
	created_at timestamptz NOT NULL DEFAULT now(),
	updated_at timestamptz NULL,
	CONSTRAINT events_pk PRIMARY KEY (id)
);
CREATE INDEX events_date_idx ON events USING btree (date);
CREATE UNIQUE INDEX events_user_id_idx ON events USING btree (user_id, date);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE events;
-- +goose StatementEnd
