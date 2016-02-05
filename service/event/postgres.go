package event

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/tapglue/multiverse/platform/metrics"
)

const (
	pgActiveByPeriod = `SELECT (json_data ->> 'user_id')::BIGINT AS userid, COUNT(*)
    FROM %s.events
      WHERE
        %s
    GROUP BY userid
    ORDER BY COUNT DESC`
	pgClauseByDay   = `(json_data ->> 'updated_at')::DATE > current_date - interval '1 day'`
	pgClauseByWeek  = `(json_data ->> 'updated_at')::DATE > current_date - interval '1 week'`
	pgClauseByMonth = `(json_data ->> 'updated_at')::DATE > current_date - interval '1 month'`

	pgCreatedByDay = `SELECT count(*), to_date(json_data->>'created_at', 'YYYY-MM-DD') as bucket
		FROM %s.events
		WHERE (json_data->>'created_at')::DATE >= '%s'
		AND (json_data->>'created_at')::DATE <= '%s'
		GROUP BY bucket
		ORDER BY bucket`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.events
		(json_data JSONB NOT NULL)`
	pgDropTable = `DROP TABLE IF EXISTS %s.objects`
)

type pgService struct {
	db *sqlx.DB
}

// NewPostgresService returns a Postgres based Service implementation.
func NewPostgresService(db *sqlx.DB) Service {
	return &pgService{db: db}
}

func (s *pgService) ActiveUserIDs(
	ns string,
	p Period,
) ([]uint64, error) {
	var clause string

	switch p {
	case ByDay:
		clause = pgClauseByDay
	case ByWeek:
		clause = pgClauseByWeek
	case ByMonth:
		clause = pgClauseByMonth
	default:
		return nil, fmt.Errorf("period %s not supported", p)
	}

	query := fmt.Sprintf(pgActiveByPeriod, ns, clause)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ids := []uint64{}
	for rows.Next() {
		var (
			id    uint64
			count int
		)

		err := rows.Scan(&id, &count)
		if err != nil {
			return nil, err
		}

		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ids, nil
}

func (s *pgService) CreatedByDay(
	ns string,
	start, end time.Time,
) (metrics.Timeseries, error) {
	query := fmt.Sprintf(
		pgCreatedByDay,
		ns,
		start.Format(metrics.BucketFormat),
		end.Format(metrics.BucketFormat),
	)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ts := []metrics.Datapoint{}
	for rows.Next() {
		var (
			bucket time.Time
			value  int
		)

		err := rows.Scan(&value, &bucket)
		if err != nil {
			return nil, err
		}

		ts = append(
			ts,
			metrics.Datapoint{
				Bucket: bucket.Format(metrics.BucketFormat),
				Value:  value,
			},
		)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ts, nil
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		wrapNamespace(pgCreateSchema, ns),
		wrapNamespace(pgCreateTable, ns),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) Teardown(ns string) error {
	qs := []string{
		wrapNamespace(pgDropTable, ns),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func wrapNamespace(query, namespace string) string {
	return fmt.Sprintf(query, namespace)
}
