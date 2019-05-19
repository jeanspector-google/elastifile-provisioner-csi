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

import (
	"time"
)

type Projects struct {

	// Project numeric id, which is automatically assigned when you create Google cloud project.
	Id string `json:"id"`

	// Project ID, which is a unique identifier for the project.
	Name string `json:"name"`

	// Project display name.
	DisplayName string `json:"displayName"`

	// List of users allowed to access resources on the specified project.
	AllowedUsers []AllowedUser `json:"allowedUsers,omitempty"`

	// Alpha features enabled on this project
	AlphaEnabled bool `json:"alphaEnabled,omitempty"`

	// The status of the project, which can be one of the following: PENDING_APPROVAL, ENABLED, or DISABLED.
	Status string `json:"status,omitempty"`

	// [Output Only] Creation timestamp in RFC3339 text format.
	CreationTimestamp time.Time `json:"creationTimestamp,omitempty"`
}
