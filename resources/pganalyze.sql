CREATE SCHEMA pganalyze;

CREATE EXTENSION IF NOT EXISTS pg_stat_statements;

CREATE OR REPLACE FUNCTION pganalyze.get_stat_statements() RETURNS SETOF pg_stat_statements AS
$$
  SELECT * FROM public.pg_stat_statements
  WHERE dbid IN (SELECT oid FROM pg_database WHERE datname = current_database());
$$ LANGUAGE sql VOLATILE SECURITY DEFINER;

CREATE OR REPLACE FUNCTION pganalyze.get_stat_activity() RETURNS SETOF pg_stat_activity AS
$$
  SELECT * FROM pg_catalog.pg_stat_activity
  WHERE datname = current_database();
$$ LANGUAGE sql VOLATILE SECURITY DEFINER;

CREATE USER pganalyze PASSWORD 'QMAfucbWKzYBhmRJaepafKzP';
REVOKE ALL ON SCHEMA public FROM pganalyze;
GRANT USAGE ON SCHEMA pganalyze TO pganalyze;
GRANT USAGE ON SCHEMA tiger TO pganalyze;