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

type Operation struct {

	// [Output Only] The unique identifier for the resource. This identifier is defined by the server.
	Id string `json:"id,omitempty"`

	// [Output Only] Name of the resource.
	Name string `json:"name,omitempty"`

	// [Output Only] A textual description of the resource.
	Description string `json:"description,omitempty"`

	// [Output Only] The value of requestId if you provided it in the request. Not present otherwise.
	ClientOperationId string `json:"clientOperationId,omitempty"`

	// [Output Only] The type of operation, such as insert, update, or delete, and so on.
	OperationType string `json:"operationType,omitempty"`

	// [Output Only] The URL of the resource that the operation modifies.
	TargetLink string `json:"targetLink,omitempty"`

	// [Output Only] The unique target ID
	TargetId string `json:"targetId,omitempty"`

	// [Output Only] The status of the operation, which can be one of the following: PENDING, RUNNING, or DONE.
	Status string `json:"status,omitempty"`

	// [Output Only] An optional textual description of the current status of the operation.
	StatusMessage string `json:"statusMessage,omitempty"`

	// [Output Only] User who requested the operation, for example: user@example.com
	User string `json:"user,omitempty"`

	// [Output Only] An optional progress indicator that ranges from 0 to 100. There is no requirement that this be linear or support any granularity of operations. This should not be used to guess when the operation will be complete. This number should monotonically increase as the operation progresses.
	Progress int32 `json:"progress,omitempty"`

	// [Output Only] The time that this operation was requested. This value is in RFC3339 text format.
	InsertTime string `json:"insertTime,omitempty"`

	// [Output Only] The time that this operation was started by the server. This value is in RFC3339 text format.
	StartTime string `json:"startTime,omitempty"`

	// [Output Only] The time that this operation was completed. This value is in RFC3339 text format.
	EndTime string `json:"endTime,omitempty"`

	// [Output Only] If errors are generated during processing of the operation, this field will be populated.
	Error_ ModelError `json:"error,omitempty"`

	// [Output Only] If warning messages are generated during processing of the operation, this field will be populated.
	Warnings []Warnings `json:"warnings,omitempty"`

	// [Output Only] This field contains the HTTP error status code that was returned. For example, a 404 means the resource was not found.
	HttpErrorStatusCode int32 `json:"httpErrorStatusCode,omitempty"`

	// [Output Only] This field contains the HTTP error message that was returned, such as NOT FOUND.
	HttpErrorMessage string `json:"httpErrorMessage,omitempty"`
}
