package aerospike

import (
	as "github.com/aerospike/aerospike-client-go"
	"github.com/tapglue/backend/config"
)

// TODO maybe extract this into aerospike/client and make it an object?

// TODO find a better way to store the buckets / sets (and maybe actually learn what those things are in the first place?)

var aerospikeClient *as.Client

// InitAerospike initializes the Aerospike client
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

// operate
func operate(namespace, set string, key interface{}, ops []*as.Operation, policy *as.WritePolicy) (rec *as.Record, err error) {
	var asKey *as.Key
	if asKey, err = as.NewKey(namespace, set, key); err != nil {
		return nil, err
	}

	return aerospikeClient.Operate(policy, asKey, ops...)
}

// Operate runs a series of operations
func OperateString(namespace, set, key string, ops []*as.Operation, policy *as.WritePolicy) (rec *as.Record, err error) {
	return operate(namespace, set, key, ops, policy)
}

// put
func put(namespace, set string, key interface{}, bins []*as.Bin, policy *as.WritePolicy) (err error) {
	var asKey *as.Key
	asKey, err = as.NewKey(namespace, set, key)
	if err != nil {
		return
	}

	if err = aerospikeClient.PutBins(policy, asKey, bins...); err != nil {
		return
	}

	return
}

// PutInt64 will put the bins into the namespace / set / key combination using the specified write policy
func PutInt64(namespace, set string, key int64, bins []*as.Bin, policy *as.WritePolicy) (err error) {
	return put(namespace, set, key, bins, policy)
}

// PutInt will put the bins into the namespace / set / key combination using the specified write policy
func PutInt(namespace, set string, key int, bins []*as.Bin, policy *as.WritePolicy) (err error) {
	return put(namespace, set, key, bins, policy)
}

// PutString will put the bins into the namespace / set / key combination using the specified write policy
func PutString(namespace, set, key string, bins []*as.Bin, policy *as.WritePolicy) (err error) {
	return put(namespace, set, key, bins, policy)
}

// get
func get(namespace, set string, key interface{}, policy *as.BasePolicy) (rec *as.Record, err error) {
	var asKey *as.Key
	if asKey, err = as.NewKey(namespace, set, key); err != nil {
		return nil, err
	}

	if policy == nil {
		policy = as.NewPolicy()
	}

	return aerospikeClient.Get(policy, asKey)
}

// GetByInt64 returns a record by searching for it via the primary key which is int64 using the specified read policy
func GetByInt64(namespace, set string, key int64, policy *as.BasePolicy) (rec *as.Record, err error) {
	return get(namespace, set, key, policy)
}

// BinToInt converts a bin value to int
func BinToInt(bin interface{}) int {
	return bin.(int)
}

// BinToInt64 converts a bin value to int64
func BinToInt64(bin interface{}) int64 {
	switch bin.(type) {
	case int:
		{
			return int64(bin.(int))
		}
	default:
		{
			return bin.(int64)
		}
	}
}

// BinToString converts a Aerospike bin string to a string
func BinToString(bin interface{}) string {
	return bin.(string)
}

// BinToBool converts a Aerospike bin boolean to a boolean
func BinToBool(bin interface{}) bool {
	return bin.(int) != 0
}

// BoolToBin converts a boolean value to a Aerospike bin boolean
func BoolToBin(b bool) int {
	if b {
		return 1
	}

	return 0
}
