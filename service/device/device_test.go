package device

import "testing"

func TestValidate(t *testing.T) {
	var (
		d  = testDevice()
		ds = List{
			{}, // Missing DeviceID
			{DeviceID: d.DeviceID},                                       // Missing EndpointARN
			{DeviceID: d.DeviceID},                                       // Missing Platform
			{DeviceID: d.DeviceID, Platform: 2},                          // Unsupported Platform
			{DeviceID: d.DeviceID, Platform: d.Platform},                 // Missing Token
			{DeviceID: d.DeviceID, Platform: d.Platform, Token: d.Token}, // Missing UserID
		}
	)

	for _, d := range ds {
		if have, want := d.Validate(), ErrInvalidDevice; !IsInvalidDevice(have) {
			t.Errorf("have %v, want %v", have, want)
		}
	}
}
