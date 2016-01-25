package object

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/tapglue/multiverse/platform/flake"
)

const (
	pgInsertObject = `INSERT INTO %s.objects(json_data) VALUES($1)`
	pgUpdateObject = `UPDATE %s.objects SET json_data = $1
		WHERE (json_data->>'id')::BIGINT = $2::BIGINT`
	pgDeleteObject = `DELETE FROM %s.objects
		WHERE (json_data->>'id')::BIGINT = $1::BIGINT`

	pgSelectObjects = `SELECT json_data FROM %s.objects
		WHERE (json_data->>'deleted')::BOOL = ?::BOOL`
	pgClauseExternalID = `AND (json_data->>'external_id')::TEXT IN (?)`
	pgClauseID         = `AND (json_data->>'id')::BIGINT = ?::BIGINT`
	pgClauseObjectID   = `AND (json_data->>'object_id')::BIGINT IN (?)`
	pgClauseOwnerID    = `AND (json_data->>'owner_id')::BIGINT IN (?)`
	pgClauseOwned      = `AND (json_data->>'owned')::BOOL = ?::BOOL`
	pgClauseType       = `AND (json_data->>'type')::TEXT IN (?)`
	pgClauseVisibility = `AND (json_data->>'visibility')::INT IN (?)`
	pgOrderCreatedAt   = `ORDER BY json_data->>'created_at' DESC LIMIT 200`

	pgCreateSchema = `CREATE SCHEMA IF NOT EXISTS %s`
	pgCreateTable  = `CREATE TABLE IF NOT EXISTS %s.objects
		(json_data JSONB NOT NULL)`

	// To ensure idempotence we want to create the index only if it doesn't exist,
	// while this feature is about to hit Postgres in 9.5 it is not yet available.
	// We fallback to a conditional create taken from:
	// http://dba.stackexchange.com/a/35626.
	pgCreateIndexCreatedAt = `CREATE INDEX object_created_at ON %s.objects
		USING btree ((json_data ->> 'created_at'))`
	pgCreateIndexExternalID = `CREATE INDEX object_external_id ON %s.objects
		USING btree (((json_data->>'external_id')::TEXT))`
	pgCreateIndexID = `CREATE INDEX object_id ON %s.objects
		USING btree (((json_data ->> 'id') :: BIGINT))`
	pgCreateIndexObjectID = `CREATE INDEX object_object_id ON %s.objects
		USING btree (((json_data->>'object_id') :: BIGINT))`
	pgCreateIndexOwnerID = `CREATE INDEX object_owner_id ON %s.objects
		USING btree (((json_data ->> 'owner_id') :: BIGINT))`
	pgCreateIndexOwned = `CREATE INDEX object_owned ON %s.objects
		USING btree (((json_data ->> 'owned') :: BOOL))`
	pgCreateIndexType = `CREATE INDEX obect_type ON %s.objects
		USING btree (((json_data->>'type') :: TEXT))`
	pgCreateIndexVisibility = `CREATE INDEX object_visibility ON %s.objects
		USING btree (((json_data ->> 'visibility') :: INT))`
	pgGuardIndex = `DO $$
		BEGIN
		IF NOT EXISTS (
			SELECT 1 FROM pg_indexes WHERE schemaname = '%s' AND indexname = '%s'
		) THEN
		%s;
		END IF;
		END$$;`

	pgDropTable = `DROP TABLE IF EXISTS %s.objects`
)

type pgService struct {
	db *sqlx.DB
}

// NewPostgresService returns a Postgres based Service implementation.
func NewPostgresService(db *sqlx.DB) Service {
	return &pgService{
		db: db,
	}
}

func (s *pgService) Put(ns string, object *Object) (*Object, error) {
	var (
		query = pgUpdateObject

		params []interface{}
	)

	if err := object.Validate(); err != nil {
		return nil, err
	}

	if object.ObjectID != 0 {
		os, err := s.Query(ns, QueryOptions{
			ID: &object.ObjectID,
		})
		if err != nil {
			return nil, err
		}

		if len(os) != 1 {
			return nil, ErrMissingReference
		}
	}

	if object.ID != 0 {
		params = []interface{}{
			object.ID,
		}

		os, err := s.Query(ns, QueryOptions{
			ID: &object.ID,
		})
		if err != nil {
			return nil, err
		}

		if len(os) == 0 {
			return nil, ErrNotFound
		}

		object.CreatedAt = os[0].CreatedAt
	} else {
		id, err := flake.NextID(flakeNamespace(ns))
		if err != nil {
			return nil, err
		}

		object.CreatedAt = time.Now().UTC()
		object.ID = id

		query = pgInsertObject
	}

	object.UpdatedAt = time.Now().UTC()

	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	params = append([]interface{}{data}, params...)

	_, err = s.db.Exec(wrapNamespace(query, ns), params...)
	if err != nil {
		// FIXME(xla): This is a defensive measure until we have proper Setup of
		// namespaces for all dependent services of an app.
		if pgWrapError(err) == ErrNamespaceNotFound {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}
			if _, err := s.db.Exec(wrapNamespace(query, ns), params...); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return object, nil
}

func (s *pgService) Query(ns string, opts QueryOptions) ([]*Object, error) {
	var (
		clauses = []string{
			wrapNamespace(pgSelectObjects, ns),
		}
		params = []interface{}{
			opts.Deleted,
		}
	)

	if len(opts.ExternalIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.ExternalIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseExternalID, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if opts.ID != nil {
		params = append(params, *opts.ID)
		clauses = append(clauses, pgClauseID)
	}

	if len(opts.OwnerIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.OwnerIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseOwnerID, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.ObjectIDs) > 0 {
		ps := []interface{}{}

		for _, id := range opts.ObjectIDs {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseObjectID, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if opts.Owned != nil {
		clause, _, err := sqlx.In(pgClauseOwned, []interface{}{*opts.Owned})
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, *opts.Owned)
	}

	if len(opts.Types) > 0 {
		ps := []interface{}{}

		for _, id := range opts.Types {
			ps = append(ps, id)
		}

		clause, _, err := sqlx.In(pgClauseType, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	if len(opts.Visibilities) > 0 {
		ps := []interface{}{}

		for _, v := range opts.Visibilities {
			ps = append(ps, v)
		}

		clause, _, err := sqlx.In(pgClauseVisibility, ps)
		if err != nil {
			return nil, err
		}

		clauses = append(clauses, clause)
		params = append(params, ps...)
	}

	clauses = append(clauses, pgOrderCreatedAt)

	os, err := s.queryObjects(strings.Join(clauses, "\n"), params...)
	if err != nil {
		// FIXME(xla): This is a defensive measure until we have proper Setup of
		// namespaces for all dependent services of an app.
		if pgWrapError(err) == ErrNamespaceNotFound {
			if err := s.Setup(ns); err != nil {
				return nil, err
			}

			os, err := s.queryObjects(strings.Join(clauses, "\n"), params...)
			if err != nil {
				return nil, err
			}

			return os, nil
		}

		return nil, err
	}

	return os, nil
}

// Remove issues a hard delete of the object with the id given.
func (s *pgService) Remove(ns string, id uint64) error {
	_, err := s.db.Exec(wrapNamespace(pgDeleteObject, ns), id)
	return pgWrapError(err)
}

func (s *pgService) Setup(ns string) error {
	qs := []string{
		wrapNamespace(pgCreateSchema, ns),
		wrapNamespace(pgCreateTable, ns),
		pgGuard("object_created_at", ns, pgCreateIndexCreatedAt),
		pgGuard("object_external_id", ns, pgCreateIndexExternalID),
		pgGuard("object_id", ns, pgCreateIndexID),
		pgGuard("object_object_id", ns, pgCreateIndexObjectID),
		pgGuard("object_owned", ns, pgCreateIndexOwned),
		pgGuard("object_owner_id", ns, pgCreateIndexOwnerID),
		pgGuard("object_type", ns, pgCreateIndexType),
		pgGuard("object_visibility", ns, pgCreateIndexVisibility),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) Teardown(namespace string) error {
	qs := []string{
		fmt.Sprintf(pgDropTable, namespace),
	}

	for _, query := range qs {
		_, err := s.db.Exec(query)
		if err != nil {
			return fmt.Errorf("query (%s): %s", query, err)
		}
	}

	return nil
}

func (s *pgService) guardWithoutSetup(
	err error,
	object *Object,
	ns, query string,
	params ...interface{},
) (*Object, error) {
	err = pgWrapError(err)
	if err == ErrNamespaceNotFound {
		err := s.Setup(ns)
		if err != nil {
			return nil, err
		}

		_, err = s.db.Exec(wrapNamespace(query, ns), params...)
		if err != nil {
			return nil, err
		}

		return object, err
	}

	return nil, err
}

func (s *pgService) queryObjects(
	query string,
	params ...interface{},
) ([]*Object, error) {
	query = sqlx.Rebind(sqlx.DOLLAR, query)

	rows, err := s.db.Query(query, params...)
	if err != nil {
		return nil, pgWrapError(err)
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	os := []*Object{}
	for rows.Next() {
		var (
			object = &Object{}

			raw []byte
		)

		err := rows.Scan(&raw)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, object)
		if err != nil {
			return nil, err
		}

		os = append(os, object)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return os, nil
}

func pgGuard(name, namespace, query string) string {
	return fmt.Sprintf(
		pgGuardIndex,
		namespace,
		name,
		wrapNamespace(query, namespace),
	)
}

func pgWrapError(err error) error {
	if err, ok := err.(*pq.Error); ok && err.Code == "42P01" {
		return ErrNamespaceNotFound
	}

	return err
}

func wrapNamespace(query, namespace string) string {
	return fmt.Sprintf(query, namespace)
}
