package controller

import (
	"github.com/tapglue/multiverse/platform/generate"
	"github.com/tapglue/multiverse/service/device"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

var defaultDeleted = false

// DeviceDeleteFunc removes the device of a user.
type DeviceDeleteFunc func(*v04_entity.Application, Origin, string) error

// DeviceDelete removes the device of a user.
func DeviceDelete(devices device.Service) DeviceDeleteFunc {
	return func(
		app *v04_entity.Application,
		origin Origin,
		deviceID string,
	) error {
		ds, err := devices.Query(app.Namespace(), device.QueryOptions{
			Deleted: &defaultDeleted,
			DeviceIDs: []string{
				deviceID,
			},
			UserIDs: []uint64{
				origin.UserID,
			},
		})
		if err != nil {
			return err
		}

		if len(ds) == 0 {
			return nil
		}

		d := ds[0]
		d.Deleted = true

		_, err = devices.Put(app.Namespace(), d)

		return err
	}
}

// DeviceUpdateFunc stores the device data and updates the endpoint.
type DeviceUpdateFunc func(
	app *v04_entity.Application,
	origin Origin,
	deviceID string,
	platform device.Platform,
	token string,
) error

// DeviceUpdate stores the device info in the given device service.
func DeviceUpdate(devices device.Service) DeviceUpdateFunc {
	return func(
		app *v04_entity.Application,
		origin Origin,
		deviceID string,
		platform device.Platform,
		token string,
	) error {
		ds, err := devices.Query(app.Namespace(), device.QueryOptions{
			Deleted: &defaultDeleted,
			DeviceIDs: []string{
				deviceID,
			},
			UserIDs: []uint64{
				origin.UserID,
			},
		})
		if err != nil {
			return err
		}

		if len(ds) > 0 && ds[0].Token == token {
			return nil
		}

		var d *device.Device

		if len(ds) > 0 {
			d = ds[0]
			d.Token = token

			// TODO: Update Endpoint.
		} else {
			// TODO: Create Endpoint.

			d = &device.Device{
				DeviceID:    deviceID,
				EndpointARN: generate.RandomString(18),
				Platform:    platform,
				Token:       token,
				UserID:      origin.UserID,
			}
		}

		_, err = devices.Put(app.Namespace(), d)
		if err != nil {
			if device.IsInvalidDevice(err) {
				return ErrInvalidEntity
			}
		}

		return err
	}
}
