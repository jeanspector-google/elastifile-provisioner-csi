/*
 * Elastifile FaaS API
 *
 * Elastifile Filesystem as a Service API
 *
 * OpenAPI spec version: 1.0
 *
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package efaasapi

type Statistics struct {

	// Number of reads in bytes
	ReadThroughput int32 `json:"readThroughput,omitempty"`

	// Number of writes in bytes
	WriteThroughput int32 `json:"writeThroughput,omitempty"`

	// Number of read IOPS
	ReadIOPS int32 `json:"readIOPS,omitempty"`

	// Number of write IOPS
	WriteIOPS int32 `json:"writeIOPS,omitempty"`

	// Number of metadata IOPS
	MdIOPS int32 `json:"mdIOPS,omitempty"`

	// Read latency in nano seconds
	ReadLatency int32 `json:"readLatency,omitempty"`

	// Write latency in nano seconds
	WriteLatency int32 `json:"writeLatency,omitempty"`

	// Metadata latency in nano seconds
	MdLatency int32 `json:"mdLatency,omitempty"`
}
