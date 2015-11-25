-- SELECT # users, connections, events, users which retrieved a feed all time
SELECT
  (SELECT count(*)
   FROM app_309_443.users)  AS "users",
  (SELECT count(*)
   FROM app_309_443.sessions)  AS "sessions",
  (SELECT count(*)
   FROM app_309_443.connections)  AS "connections",
  (SELECT count(*)
   FROM app_309_443.events) AS "events",
  (SELECT count(*)
   FROM app_309_443.users
   WHERE last_read <> '2015-05-01 01:23:45.000000') AS "users which retrieved a feed";

-- SELECT any activity in the last 30 days
SELECT
  (SELECT count(*)
   FROM app_309_443.users
   WHERE (json_data->>'updated_at')::TIMESTAMP > (current_date - interval '30 days'))  AS "users",
  (SELECT count(*)
   FROM app_309_443.sessions
   WHERE (created_at)::TIMESTAMP > (current_date - interval '30 days'))  AS "sessions",
  (SELECT count(*)
   FROM app_309_443.connections
   WHERE (json_data->>'updated_at')::TIMESTAMP > (current_date - interval '30 days'))  AS "connections",
  (SELECT count(*)
   FROM app_309_443.events
   WHERE (json_data->>'updated_at')::TIMESTAMP > (current_date - interval '30 days')) AS "events",
  (SELECT count(*)
   FROM app_309_443.users
   WHERE (last_read)::TIMESTAMP > (current_date - interval '30 days')) AS "users which retrieved a feed";

-- SELECT creation activity in the last 30 days
SELECT
  (SELECT count(*)
   FROM app_309_443.users
   WHERE (json_data->>'created_at')::TIMESTAMP > (current_date - interval '30 days'))  AS "users",
  (SELECT count(*)
   FROM app_309_443.sessions
   WHERE (created_at)::TIMESTAMP > (current_date - interval '30 days'))  AS "sessions",
  (SELECT count(*)
   FROM app_309_443.connections
   WHERE (json_data->>'created_at')::TIMESTAMP > (current_date - interval '30 days'))  AS "connections",
  (SELECT count(*)
   FROM app_309_443.events
   WHERE (json_data->>'created_at')::TIMESTAMP > (current_date - interval '30 days')) AS "events",
  (SELECT count(*)
   FROM app_309_443.users
   WHERE (last_read)::TIMESTAMP > (current_date - interval '30 days')) AS "users which retrieved a feed";


-- Users updated, created after the initial import or fetched a feed
SELECT count(*)
FROM app_309_443.users
WHERE
  ((json_data ->> 'updated_at') :: DATE > '2015-10-14 10:00:00' OR
   (json_data ->> 'created_at') :: DATE > '2015-10-14 10:00:00')
  OR last_read <> '2015-05-01 01:23:45.000000';

-- New users after the initial import
SELECT count(*)
FROM app_309_443.users
WHERE (json_data ->> 'created_at') :: DATE > '2015-10-14 10:00:00';

-- Users who created an event
SELECT count(DISTINCT (json_data ->> 'user_id'))
FROM app_309_443.events;

-- Average events per user who created an event
WITH generated_events(user_id, event_count) AS (
    SELECT
      app_309_443.users.json_data ->> 'id' AS user_id,
      count(events.json_data ->> 'id')     AS event_count
    FROM app_309_443.events
      JOIN app_309_443.users ON app_309_443.events.json_data ->> 'user_id' = app_309_443.users.json_data ->> 'id'
    GROUP BY app_309_443.users.json_data ->> 'id'
)
SELECT avg(event_count)
FROM generated_events;

-- Sum events non-tg_follow
WITH generated_events(user_id, event_count) AS (
    SELECT
      app_309_443.users.json_data ->> 'id' AS user_id,
      count(events.json_data ->> 'id')     AS event_count
    FROM app_309_443.events
      JOIN app_309_443.users ON app_309_443.events.json_data ->> 'user_id' = app_309_443.users.json_data ->> 'id'
    WHERE (app_309_443.events.json_data->>'type')::TEXT <> 'tg_follow'
    GROUP BY app_309_443.users.json_data ->> 'id'
)
SELECT sum(event_count)
FROM generated_events;

-- Sum events tg_follow
WITH generated_events(user_id, event_count) AS (
    SELECT
      app_309_443.users.json_data ->> 'id' AS user_id,
      count(events.json_data ->> 'id')     AS event_count
    FROM app_309_443.events
      JOIN app_309_443.users ON app_309_443.events.json_data ->> 'user_id' = app_309_443.users.json_data ->> 'id'
    WHERE (app_309_443.events.json_data->>'type')::TEXT = 'tg_follow'
    GROUP BY app_309_443.users.json_data ->> 'id'
)
SELECT sum(event_count)
FROM generated_events;

-- Sum events non-tg_follow after import date
WITH generated_events(user_id, event_count) AS (
    SELECT
      app_309_443.users.json_data ->> 'id' AS user_id,
      count(events.json_data ->> 'id')     AS event_count
    FROM app_309_443.events
      JOIN app_309_443.users ON app_309_443.events.json_data ->> 'user_id' = app_309_443.users.json_data ->> 'id'
    WHERE (app_309_443.events.json_data->>'type')::TEXT <> 'tg_follow'
          AND (app_309_443.events.json_data ->> 'created_at') :: DATE > '2015-10-14 10:00:00'
    GROUP BY app_309_443.users.json_data ->> 'id'
)
SELECT sum(event_count)
FROM generated_events;

-- Sum events tg_follow after import date
WITH generated_events(user_id, event_count) AS (
    SELECT
      app_309_443.users.json_data ->> 'id' AS user_id,
      count(events.json_data ->> 'id')     AS event_count
    FROM app_309_443.events
      JOIN app_309_443.users ON app_309_443.events.json_data ->> 'user_id' = app_309_443.users.json_data ->> 'id'
    WHERE (app_309_443.events.json_data->>'type')::TEXT = 'tg_follow'
      AND (app_309_443.events.json_data ->> 'created_at') :: DATE > '2015-10-14 10:00:00'
    GROUP BY app_309_443.users.json_data ->> 'id'
)
SELECT sum(event_count)
FROM generated_events;