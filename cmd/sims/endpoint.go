package main

const (
	serviceSNS          = "SNS"
	typeDeliveryFailure = "DeliveryFailure"
)

func updateEndpoint(c endpointChange, disableDevice disableDeviceFunc) (err error) {
	defer func() {
		if err == nil {
			c.ack()
		}
	}()

	if c.Service != serviceSNS {
		return nil
	}

	if c.EventType != typeDeliveryFailure {
		return nil
	}

	return disableDevice(c.Resource, c.EndpointArn)
}
