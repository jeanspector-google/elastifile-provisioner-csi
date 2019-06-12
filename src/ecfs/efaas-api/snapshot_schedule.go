/* 
 * Elastifile FaaS API
 *
 * Elastifile Filesystem as a Service API
 *
 * OpenAPI spec version: 2.0
 * 
 * Generated by: https://github.com/swagger-api/swagger-codegen.git
 */

package EfaasApi

type SnapshotSchedule struct {

	// Indicates whether snapshots usage in enabled for the instance.
	Enable bool `json:"enable"`

	// When snapshot.enable is set to true, this field indicates how often to schedule snapshot creation.
	Schedule string `json:"schedule"`

	// Snapshot retention policy. The number of days to hold the snapshot till automatic deletion. Default 0, meaning no retention policy defined.
	Retention float32 `json:"retention,omitempty"`
}