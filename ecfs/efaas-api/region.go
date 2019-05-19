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

type Region struct {

	// [Output Only] The unique identifier for the resource. This identifier is defined by the server.
	Id string `json:"id,omitempty"`

	// [Output Only] Name of the resource.
	Name string `json:"name"`

	// Geographical location in which the region resides in.
	Location string `json:"location"`

	// [Output Only] A list of zones available in this region
	Zones []string `json:"zones,omitempty"`
}
