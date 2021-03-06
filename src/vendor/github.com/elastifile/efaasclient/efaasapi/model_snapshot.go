/*
 * Elastifile FaaS API
 *
 * Elastifile Filesystem as a Service API
 *
 * API version: 2.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package efaasapi

type Snapshot struct {
	// The name of the resource, provided by the client when initially creating the resource. The resource name must be 1-63 characters long, and comply with RFC1035
	Name string `json:"name"`
	// Snapshot retention policy. The number of days to hold the snapshot till automatic deletion. Default 0, meaning no retention policy defined.
	Retention float32 `json:"retention"`
}
