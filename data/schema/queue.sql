CREATE TYPE job_status AS ENUM ('claimed', 'launched');

CREATE TABLE IF NOT EXISTS jobs (
  job_id      bigserial NOT NULL PRIMARY KEY,
  next_try    bigint NOT NULL CHECK (next_try > 0),
  grabber_id  bigint,
  status      job_status NOT NULL DEFAULT 'claimed',
  created     timestamp  NOT NULL DEFAULT CURRENT_TIMESTAMP,
  retry_count integer NOT NULL DEFAULT 0 CHECK (retry_count >= 0),
  retry_delay integer NOT NULL DEFAULT 0 CHECK (retry_delay >= 0),
  fail_count  integer NOT NULL DEFAULT 0 CHECK (fail_count >= 0),
  name        varchar(255)  NOT NULL DEFAULT '',
  url         varchar(512)  NOT NULL DEFAULT '',
  payload     json NOT NULL DEFAULT '[]'::json,
  timeout     integer NOT NULL DEFAULT 0 CONSTRAINT positive_timeout CHECK (timeout >= 0)
);

CREATE INDEX IF NOT EXISTS jobs_status_idx ON jobs USING btree (status);
