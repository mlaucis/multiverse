package aerospike

import (
	as "github.com/aerospike/aerospike-client-go"
	"github.com/tapglue/backend/config"
)

var aerospikeClient *as.Client

func InitAerospike(aerospike config.Aerospike) (err error) {
	policy := as.NewClientPolicy()
	policy.Timeout = aerospike.ConnectionTimeout()
	policy.ConnectionQueueSize = aerospike.ConnectionQueueSize()
	policy.FailIfNotConnected = aerospike.ConnectOrFail()

	hosts := []*as.Host{}

	for host, port := range aerospike.Servers() {
		hosts = append(hosts, as.NewHost(host, port))
	}

	aerospikeClient, err = as.NewClientWithPolicyAndHost(policy, hosts...)

	return
}

func Client() *as.Client {
	return aerospikeClient
}
