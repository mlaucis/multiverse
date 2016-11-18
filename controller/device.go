package controller

import (
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/device"
)

var defaultDeleted = false

// DeviceDeleteFunc removes the device of a user.
type DeviceDeleteFunc func(*app.App, Origin, string) error

// DeviceDelete removes the device of a user.
func DeviceDelete(devices device.Service) DeviceDeleteFunc {
	return func(
		currentApp *app.App,
		origin Origin,
		deviceID string,
	) error {
		ds, err := devices.Query(currentApp.Namespace(), device.QueryOptions{
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

		_, err = devices.Put(currentApp.Namespace(), d)

		return err
	}
}

// DeviceUpdateFunc stores the device data and updates the endpoint.
type DeviceUpdateFunc func(
	currentApp *app.App,
	origin Origin,
	deviceID string,
	platform device.Platform,
	token string,
	language string,
) error

// DeviceUpdate stores the device info in the given device service.
func DeviceUpdate(devices device.Service) DeviceUpdateFunc {
	return func(
		currentApp *app.App,
		origin Origin,
		deviceID string,
		platform device.Platform,
		token string,
		language string,
	) error {
		ds, err := devices.Query(currentApp.Namespace(), device.QueryOptions{
			Deleted: &defaultDeleted,
			Platforms: []device.Platform{
				platform,
			},
			UserIDs: []uint64{
				origin.UserID,
			},
		})
		if err != nil {
			return err
		}

		d := &device.Device{}

		for _, dev := range ds {
			if dev.DeviceID == deviceID && dev.Token == token {
				return nil
			}

			if dev.DeviceID == deviceID || dev.Token == token {
				d = dev
			}
		}

		d.DeviceID = deviceID
		d.Disabled = false
		d.Language = language
		d.Platform = platform
		d.Token = token
		d.UserID = origin.UserID

		_, err = devices.Put(currentApp.Namespace(), d)
		if err != nil {
			if device.IsInvalidDevice(err) {
				return wrapError(ErrInvalidEntity, "%s", err)
			}
		}

		return err
	}
}
