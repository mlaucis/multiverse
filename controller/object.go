package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// ObjectController bundles business contraints for objects.
type ObjectController struct {
	connections connection.Service
	objects     object.Service
}

// NewObjectController returns a controller instance.
func NewObjectController(
	connections connection.Service,
	objects object.Service,
) *ObjectController {
	return &ObjectController{
		connections: connections,
		objects:     objects,
	}
}

// Create associates the given Object with the owner and stores it in the
// service.
func (c *ObjectController) Create(
	app *v04_entity.Application,
	o *object.Object,
	origin uint64,
) (*object.Object, error) {
	o.OwnerID = origin

	return c.objects.Put(app.Namespace(), o)
}

// Delete marks an Object as deleted and updates the service.
func (c *ObjectController) Delete(app *v04_entity.Application, id uint64) error {
	var o *object.Object

	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &id,
	})
	if err != nil {
		return err
	}

	// A delete should be idempotent and always succeed.
	if len(os) == 0 {
		return nil
	}

	o = os[0]
	o.Deleted = true

	_, err = c.objects.Put(app.Namespace(), o)
	if err != nil {
		return err
	}

	return nil
}

// List returns the list of objects for the given user id.
func (c *ObjectController) List(app *v04_entity.Application, ownerID uint64) ([]*object.Object, error) {
	return c.objects.Query(app.Namespace(), object.QueryOptions{
		OwnerIDs: []uint64{
			ownerID,
		},
	})
}

// ListAll returns the list of objects for the given app.
func (c *ObjectController) ListAll(app *v04_entity.Application) ([]*object.Object, error) {
	return c.objects.Query(app.Namespace(), object.QueryOptions{
		Visibilities: []object.Visibility{
			object.VisibilityPublic,
			object.VisibilityGlobal,
		},
	})
}

// ListConnections returns the list of objects for the connections of the given
// user.
func (c *ObjectController) ListConnections(
	app *v04_entity.Application,
	originID uint64,
) ([]*object.Object, error) {
	cs, err := c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{
			originID,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		Types: []connection.Type{
			connection.TypeFriend,
			connection.TypeFollow,
		},
	})
	if err != nil {
		return nil, err
	}

	ids := cs.ToIDs()

	cs, err = c.connections.Query(app.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		ToIDs: []uint64{
			originID,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		Types: []connection.Type{
			connection.TypeFriend,
		},
	})
	if err != nil {
		return nil, err
	}

	ids = append(ids, cs.FromIDs()...)

	if len(ids) == 0 {
		return []*object.Object{}, nil
	}

	return c.objects.Query(app.Namespace(), object.QueryOptions{
		OwnerIDs: ids,
		Visibilities: []object.Visibility{
			object.VisibilityConnection,
			object.VisibilityPublic,
			object.VisibilityGlobal,
		},
	})
}

// Retrieve return the Object for the given id.
func (c *ObjectController) Retrieve(
	app *v04_entity.Application,
	id uint64,
) (*object.Object, error) {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	if len(os) != 1 {
		return nil, ErrNotFound
	}

	return os[0], nil
}

// Update stores the new object in the service.
func (c *ObjectController) Update(
	app *v04_entity.Application,
	id uint64,
	o *object.Object,
) (*object.Object, error) {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	if len(os) != 1 {
		return nil, ErrNotFound
	}

	// Preserve ownership infomration.
	o.ID = id
	o.OwnerID = os[0].OwnerID

	return c.objects.Put(app.Namespace(), o)
}
