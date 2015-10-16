DROP SCHEMA IF EXISTS public CASCADE ;

CREATE SCHEMA IF NOT EXISTS public;

CREATE EXTENSION postgis;
CREATE EXTENSION postgis_topology;
CREATE EXTENSION fuzzystrmatch;
CREATE EXTENSION postgis_tiger_geocoder;

CREATE SCHEMA tg;

CREATE TABLE tg.accounts (
  id SERIAL PRIMARY KEY NOT NULL,
  json_data JSONB NOT NULL
);

CREATE TABLE tg.account_users (
  id SERIAL PRIMARY KEY NOT NULL,
  account_id INT NOT NULL,
  json_data JSONB NOT NULL
);

CREATE TABLE tg.account_user_sessions (
  account_id INT NOT NULL,
  account_user_id INT NOT NULL,
  session_id CHAR(40) NOT NULL,
  created_at TIMESTAMP DEFAULT now() NOT NULL
);

CREATE TABLE tg.applications (
  id SERIAL PRIMARY KEY NOT NULL,
  account_id INT NOT NULL,
  json_data JSONB NOT NULL,
  enabled INT DEFAULT 1 NOT NULL
);

CREATE TABLE tg.consumers
(
  consumer_name TEXT NOT NULL,
  consumer_position TEXT NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

CREATE INDEX on tg.accounts USING GIN (json_data jsonb_path_ops);
CREATE INDEX on tg.account_users USING GIN (json_data jsonb_path_ops);
CREATE INDEX app_public_id ON tg.applications (((json_data ->> 'id') :: TEXT));
CREATE INDEX app_token ON tg.applications (((json_data ->> 'token') :: TEXT));
CREATE INDEX app_backend_token ON tg.applications (((json_data ->> 'backend_token') :: TEXT));

