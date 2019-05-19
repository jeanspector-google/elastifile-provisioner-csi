package efaas

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"testing"
	"time"

	efaasapi "csi-provisioner-elastifile/ecfs/efaas-api"
)

var jsonData = []byte(`
	{
	 "type": "service_account",
	 "project_id": "elastifile-gce-lab-c934",
	 "private_key_id": "5e0d188967e7f23ad77129ff4c9ab59889ccd25d",
	 "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCMBJyta1PEkd7q\nCLEYNdUBqk4Hlnw7mGXnByjao+4SOZi7mJ1NIAtYjptJ/rcPxjft+hxEba1a1DON\nUU7RuJ3eQk+kLVHdbD2D4noMw6VxJtuWnuyQ2V8v5ojv8kVvVSsbkDAQHVGKTe/8\nCEHxlekGoY0NC+KwWlUKmb7cv/B/2aD1eFsyV7ALE/YJmyFbbvtLrab+U5js04ER\nIWcE+gKlvAF7Xq9Iq6MucyjRvgPagz5RSP146HjbCPdJIz3ilcEL7idVGaZnnx/P\ncZAqYnYZAJTGBhi4fUEpAYR7KVUWIVfc9oXEKJDNwwBHnyyZMBPdYn9prs7xgrEL\ngA+WHPPZAgMBAAECggEACVNhUBee66+/hhzwFqm3NzYtnknCmoGK//k1GmLiv2oA\npzYB/BoPR2WwKByD+tP786i96zzW1/7cNCRfOI6wTRZjkY7HLhVAf6E8+c6qHUA2\nTfDl1rvzoBAdvMWJJGIqzdorqVcakDiirEmsgre2Xo+yAlVxUsehdGRLFw7dqNYv\nrINMqjE2W/SCd8jw2WmplmH+c0MvBKkving9CCNgFnvSMUGinv7y3Zvf2GpplvlC\nFdSFGGXxn1o6HbgrkovKn6EVZ8nP3JadG5evwjotEv1fcEu4vOKMq/jgvfxzscRf\ng9bfdhb3/oc+x43dsH3fR0axaImB7LKKgfu7w7vnJQKBgQDCmgAE7noPd0bt7Xg+\nrl44OgCHv3x0QY4lx0y07Yo1Bg1C72H8BCghr/5rxGUOSCGjoFYTVeLhCVIsYX+8\nxbtplxCJFAgN7lu48EyCgIpP7ppjf1a3Uh762O04BCMw0tXw22ich7d4KN5+r8L7\nOknRStrZYD89QjoUsSEYOK0wnwKBgQC4MePUNoBJEG+yhlMOpDz7mnf/F1U4gFQQ\nxD4stAEA1P/QuSgMb0snJJA3yT3dCL4W2DUxDCWOH/Wx3XnJy216+QR//8fHImCR\nYS4fjmaWlbMOKko1yeCtCLsNfA5uB5Yplrujn2o6v5BE52h3JCjW4qUqzZ6T9cBq\n0rQFacWwhwKBgBKLJDdUFjOFFTA08cFfUkEfXc+RsqVNXeNBs5CGFiZpVjgroXWn\nW7+iCqdwRoTu4K276JfdFkqFXdw2yjpNyUcNixjU3NOfBASCeXfyEbv+K54Rk0zS\nuXsD0s8ErenIHXTfI3/O+u+rTVBbJURVUJVuAZ63Ki+HMQupuVKai/5XAoGBALcp\nHSV8IKsHBhtfSR5JIT8MhoCKIjsyGOYnTrBDOrAqHkveor1iujetOx/OJI80T1oG\nGzavnnSqwTXiR2XrvO1IzDnADletjptiKGxGvSrGp6vRT8QXACzwfpjVIMA3GRI4\nClSVhBvxO7PY7N90fIvaCmX629LD0FgpN8weNu/nAoGAP4rXRr37757Q+c/qeKyU\nsmUCYeHj6w+GIkqJIhsDsj5tE8fLTyU87LF6hvscxYJCX9ZVycvhuzRBiFLkc9yo\nZUKC4SllFDw4Zl63RU7me3PnZHpomiNs0hk3fgqAME1Cx3Pn8NT6iptybSqk2kb7\nHOuPCeblZecVZU0UOPyQrWM=\n-----END PRIVATE KEY-----\n",
	 "client_email": "efaas-csi@elastifile-gce-lab-c934.iam.gserviceaccount.com",
	 "client_id": "102179953128561786237",
	 "auth_uri": "https://accounts.google.com/o/oauth2/auth",
	 "token_uri": "https://oauth2.googleapis.com/token",
	 "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
	 "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/efaas-csi%40elastifile-gce-lab-c934.iam.gserviceaccount.com"
	}
	`)

func TestREST(t *testing.T) {
	t.Parallel()

	res, err := demo1(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	t.Logf("RES: %v", string(res))
}

func TestAPI(t *testing.T) {
	client, err := GetEfaasClient(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	res, err := apiCallGet(client, InstancesURL)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	t.Logf("RES: %v", string(res))
}

func TestSwaggerLowLevelAPI(t *testing.T) {
	client, err := GetEfaasClient(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	apiConf := efaasapi.NewConfiguration()
	apiConf.BasePath = BaseURL
	apiConf.AccessToken = client.GoogleIdToken
	apiConf.Debug = true
	apiConf.DebugFile = "/tmp/api-debug.log"

	// Insecure transport
	defaultTransport := http.DefaultTransport.(*http.Transport)
	apiConf.Transport = &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		TLSClientConfig:       &tls.Config{InsecureSkipVerify: true}, // TODO: FIXME before deploying to production
	}
	apiConf.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %v", client.GoogleIdToken))

	res, err := apiConf.APIClient.CallAPI("https://bronze-eagle.gcp.elastifile.com/api/v1/regions", "GET",
		nil, apiConf.DefaultHeader, nil, nil, "", nil)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}
	t.Logf("RES: %+v", res)
}

func TestOpenAPI_CreateInstance(t *testing.T) {
	efaasConf, err := NewEfaasConf(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}

	snapshots := efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 0.0,
	}

	accessor1 := efaasapi.AccessorItems{
		SourceRange:  "10.142.0.0/20", // TODO: Detect the range via K8s OR get it at deploy time
		AccessRights: "readWrite",
	}

	accessors := efaasapi.Accessors{
		Items: []efaasapi.AccessorItems{accessor1},
	}

	filesystem := efaasapi.DataContainer{
		Name:        "dc1",                          // Filesystem name
		Description: fmt.Sprintf("Filesystem desc"), // Filesystem description
		QuotaType:   "fixed",                        // Supported values are: auto and fixed. Use auto if you have one filesystem, the size of the filesystem will be the same as the instance size. Use fixed if you have more than one filesystem, and set the filesystem size through filesystemQuota.
		HardQuota:   10 * 1024 * 1024 * 1024,        // Set the size of a filesystem if filesystemQuotaType is set to fixed. If it is set to auto, this value is ignored and quota is the instance total size.
		Snapshot:    snapshots,                      // Snapshot object
		Accessors:   accessors,                      // Defines the access rights to the File System. This is a listof access rights configured by the client for the file system.
	}

	payload := efaasapi.Instances{
		Name:                     "jean-instance1",
		Description:              "eFaaS instance description",
		ServiceClass:             "capacity-optimized-az",
		Region:                   "us-east1",
		Zone:                     "us-east1-b",
		ProvisionedCapacityUnits: 3,
		Network:                  "default",
		Filesystems:              []efaasapi.DataContainer{filesystem},
	}

	op, resp, err := instancesAPI.CreateInstance(ProjectId, payload, "")
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	if resp.StatusCode > http.StatusAccepted {
		t.Fatal("HTTP request failed", "status code", resp.StatusCode, "status", resp.Status)
	}

	t.Logf("Opration: %#v", op)
	t.Logf("Response: %#v", resp)
	t.Logf("Response payload: %v", fmt.Sprint(string(resp.Payload)))

	t.Logf("Waiting for operation id %v ...", op.Id)

	err = WaitForOperationStatusComplete(efaasConf, op.Id, time.Hour)
	if err != nil {
		t.Fatal("WaitForOperationStatusComplete failed", "err", err)
	}
}

func TestOpenAPI_GetInstance(t *testing.T) { // Works (with update of int32 to int64)
	efaasConf, err := NewEfaasConf(jsonData)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	instancesAPI := efaasapi.ProjectsprojectinstancesApi{Configuration: efaasConf}

	inst, resp, err := instancesAPI.GetInstance("test-instance--efb4feee-1", ProjectId)
	if err != nil {
		t.Fatal(fmt.Sprintf("AAAAA %v", err.Error()))
	}

	if resp.StatusCode >= 300 {
		t.Fatal("HTTP request failed", "status code", resp.StatusCode, "status", resp.Status)
	}

	t.Logf("Instance: %#v", inst)
	t.Logf("Response: %#v", resp)
	t.Logf("Response payload: %v", fmt.Sprint(string(resp.Payload)))
}
