package controller

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"

	"github.com/tapglue/multiverse/platform/generate"
	"github.com/tapglue/multiverse/service/app"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// AppCreateFunc creates an application for the current Org.
type AppCreateFunc func(
	org *v04_entity.Organization,
	name string,
	description string,
) (*app.App, error)

// AppCreate creates an application for the current Org.
func AppCreate(apps app.Service) AppCreateFunc {
	return func(
		currentOrg *v04_entity.Organization,
		name, description string,
	) (*app.App, error) {
		token, backendToken, err := generateTokens()
		if err != nil {
			return nil, err
		}

		publicID, err := generate.UUID()
		if err != nil {
			return nil, err
		}

		return apps.Put(app.NamespaceDefault, &app.App{
			BackendToken: backendToken,
			Description:  description,
			Enabled:      true,
			InProduction: false,
			Name:         name,
			OrgID:        uint64(currentOrg.ID),
			PublicID:     publicID,
			PublicOrgID:  currentOrg.PublicID,
			Token:        token,
		})
	}
}

// AppDeleteFunc disables the App and renders it unusbale.
type AppDeleteFunc func(*v04_entity.Organization, string) error

// AppDelete disables the App and renders it unusable.
func AppDelete(apps app.Service) AppDeleteFunc {
	return func(org *v04_entity.Organization, publicID string) error {
		as, err := apps.Query(app.NamespaceDefault, app.QueryOptions{
			Enabled: &defaultEnabled,
			PublicIDs: []string{
				publicID,
			},
		})
		if err != nil {
			return err
		}

		if len(as) != 1 {
			return ErrNotFound
		}

		as[0].Enabled = false

		_, err = apps.Put(app.NamespaceDefault, as[0])

		return err
	}
}

// AppListFunc returns all Apps for the current Org.
type AppListFunc func(*v04_entity.Organization, app.QueryOptions) (app.List, error)

// AppList returns all Apps for the current Org.
func AppList(apps app.Service) AppListFunc {
	return func(
		currentOrg *v04_entity.Organization,
		opts app.QueryOptions,
	) (app.List, error) {
		opts.Enabled = &defaultEnabled
		opts.OrgIDs = []uint64{
			uint64(currentOrg.ID),
		}

		return apps.Query(app.NamespaceDefault, opts)
	}
}

// AppUpdateFunc updates the values of an App..
type AppUpdateFunc func(
	currentOrg *v04_entity.Organization,
	publiID string,
	name, description string,
) (*app.App, error)

// AppUpdate updates the values of an App.
func AppUpdate(apps app.Service) AppUpdateFunc {
	return func(
		currentOrg *v04_entity.Organization,
		publiID string,
		name, description string,
	) (*app.App, error) {
		as, err := apps.Query(app.NamespaceDefault, app.QueryOptions{
			Enabled: &defaultEnabled,
			PublicIDs: []string{
				publiID,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(as) != 1 {
			return nil, ErrNotFound
		}

		a := as[0]
		a.Name = name
		a.Description = description

		return apps.Put(app.NamespaceDefault, a)
	}
}

func generateTokens() (string, string, error) {
	src := rand.NewSource(time.Now().UnixNano())

	tokenHash := md5.New()
	_, err := tokenHash.Write(generate.RandomBytes(src, 32))
	if err != nil {
		return "", "", err
	}
	token := fmt.Sprintf("%x", tokenHash.Sum(nil))

	backendHash := md5.New()
	_, err = backendHash.Write(generate.RandomBytes(src, 12))
	if err != nil {
		return "", "", err
	}

	return token, fmt.Sprintf(
		"%s%s",
		token,
		fmt.Sprintf("%x", backendHash.Sum(nil))[:12],
	), nil
}
