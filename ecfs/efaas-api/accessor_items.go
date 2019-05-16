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

type AccessorItems struct {

	// Source range to accept traffic from, These range must be, expressed in CIDR format or you also use the 'all' alias to define the range of internal addresses (10.0.0.0/8, 172.16.0.0/12 and 192.168.0.0/16)
	SourceRange string `json:"sourceRange"`

	// Define access rights for clients originated from the specified source range, supported values are: readOnly and readWrite.
	AccessRights string `json:"accessRights"`
}
