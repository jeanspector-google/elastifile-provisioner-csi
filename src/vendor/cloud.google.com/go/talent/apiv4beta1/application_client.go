// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by gapic-generator. DO NOT EDIT.

package talent

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/golang/protobuf/proto"
	gax "github.com/googleapis/gax-go/v2"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"google.golang.org/api/transport"
	talentpb "google.golang.org/genproto/googleapis/cloud/talent/v4beta1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
)

// ApplicationCallOptions contains the retry settings for each method of ApplicationClient.
type ApplicationCallOptions struct {
	CreateApplication []gax.CallOption
	GetApplication    []gax.CallOption
	UpdateApplication []gax.CallOption
	DeleteApplication []gax.CallOption
	ListApplications  []gax.CallOption
}

func defaultApplicationClientOptions() []option.ClientOption {
	return []option.ClientOption{
		option.WithEndpoint("jobs.googleapis.com:443"),
		option.WithScopes(DefaultAuthScopes()...),
	}
}

func defaultApplicationCallOptions() *ApplicationCallOptions {
	retry := map[[2]string][]gax.CallOption{
		{"default", "idempotent"}: {
			gax.WithRetry(func() gax.Retryer {
				return gax.OnCodes([]codes.Code{
					codes.Unavailable,
				}, gax.Backoff{
					Initial:    100 * time.Millisecond,
					Max:        60000 * time.Millisecond,
					Multiplier: 1.3,
				})
			}),
		},
	}
	return &ApplicationCallOptions{
		CreateApplication: retry[[2]string{"default", "non_idempotent"}],
		GetApplication:    retry[[2]string{"default", "idempotent"}],
		UpdateApplication: retry[[2]string{"default", "non_idempotent"}],
		DeleteApplication: retry[[2]string{"default", "idempotent"}],
		ListApplications:  retry[[2]string{"default", "idempotent"}],
	}
}

// ApplicationClient is a client for interacting with Cloud Talent Solution API.
//
// Methods, except Close, may be called concurrently. However, fields must not be modified concurrently with method calls.
type ApplicationClient struct {
	// The connection to the service.
	conn *grpc.ClientConn

	// The gRPC API client.
	applicationClient talentpb.ApplicationServiceClient

	// The call options for this service.
	CallOptions *ApplicationCallOptions

	// The x-goog-* metadata to be sent with each request.
	xGoogMetadata metadata.MD
}

// NewApplicationClient creates a new application service client.
//
// A service that handles application management, including CRUD and
// enumeration.
func NewApplicationClient(ctx context.Context, opts ...option.ClientOption) (*ApplicationClient, error) {
	conn, err := transport.DialGRPC(ctx, append(defaultApplicationClientOptions(), opts...)...)
	if err != nil {
		return nil, err
	}
	c := &ApplicationClient{
		conn:        conn,
		CallOptions: defaultApplicationCallOptions(),

		applicationClient: talentpb.NewApplicationServiceClient(conn),
	}
	c.setGoogleClientInfo()
	return c, nil
}

// Connection returns the client's connection to the API service.
func (c *ApplicationClient) Connection() *grpc.ClientConn {
	return c.conn
}

// Close closes the connection to the API service. The user should invoke this when
// the client is no longer required.
func (c *ApplicationClient) Close() error {
	return c.conn.Close()
}

// setGoogleClientInfo sets the name and version of the application in
// the `x-goog-api-client` header passed on each request. Intended for
// use by Google-written clients.
func (c *ApplicationClient) setGoogleClientInfo(keyval ...string) {
	kv := append([]string{"gl-go", versionGo()}, keyval...)
	kv = append(kv, "gapic", versionClient, "gax", gax.Version, "grpc", grpc.Version)
	c.xGoogMetadata = metadata.Pairs("x-goog-api-client", gax.XGoogHeader(kv...))
}

// CreateApplication creates a new application entity.
func (c *ApplicationClient) CreateApplication(ctx context.Context, req *talentpb.CreateApplicationRequest, opts ...gax.CallOption) (*talentpb.Application, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", req.GetParent()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.CreateApplication[0:len(c.CallOptions.CreateApplication):len(c.CallOptions.CreateApplication)], opts...)
	var resp *talentpb.Application
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.applicationClient.CreateApplication(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetApplication retrieves specified application.
func (c *ApplicationClient) GetApplication(ctx context.Context, req *talentpb.GetApplicationRequest, opts ...gax.CallOption) (*talentpb.Application, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.GetApplication[0:len(c.CallOptions.GetApplication):len(c.CallOptions.GetApplication)], opts...)
	var resp *talentpb.Application
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.applicationClient.GetApplication(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// UpdateApplication updates specified application.
func (c *ApplicationClient) UpdateApplication(ctx context.Context, req *talentpb.UpdateApplicationRequest, opts ...gax.CallOption) (*talentpb.Application, error) {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "application.name", req.GetApplication().GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.UpdateApplication[0:len(c.CallOptions.UpdateApplication):len(c.CallOptions.UpdateApplication)], opts...)
	var resp *talentpb.Application
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		resp, err = c.applicationClient.UpdateApplication(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteApplication deletes specified application.
func (c *ApplicationClient) DeleteApplication(ctx context.Context, req *talentpb.DeleteApplicationRequest, opts ...gax.CallOption) error {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "name", req.GetName()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.DeleteApplication[0:len(c.CallOptions.DeleteApplication):len(c.CallOptions.DeleteApplication)], opts...)
	err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
		var err error
		_, err = c.applicationClient.DeleteApplication(ctx, req, settings.GRPC...)
		return err
	}, opts...)
	return err
}

// ListApplications lists all applications associated with the profile.
func (c *ApplicationClient) ListApplications(ctx context.Context, req *talentpb.ListApplicationsRequest, opts ...gax.CallOption) *ApplicationIterator {
	md := metadata.Pairs("x-goog-request-params", fmt.Sprintf("%s=%v", "parent", req.GetParent()))
	ctx = insertMetadata(ctx, c.xGoogMetadata, md)
	opts = append(c.CallOptions.ListApplications[0:len(c.CallOptions.ListApplications):len(c.CallOptions.ListApplications)], opts...)
	it := &ApplicationIterator{}
	req = proto.Clone(req).(*talentpb.ListApplicationsRequest)
	it.InternalFetch = func(pageSize int, pageToken string) ([]*talentpb.Application, string, error) {
		var resp *talentpb.ListApplicationsResponse
		req.PageToken = pageToken
		if pageSize > math.MaxInt32 {
			req.PageSize = math.MaxInt32
		} else {
			req.PageSize = int32(pageSize)
		}
		err := gax.Invoke(ctx, func(ctx context.Context, settings gax.CallSettings) error {
			var err error
			resp, err = c.applicationClient.ListApplications(ctx, req, settings.GRPC...)
			return err
		}, opts...)
		if err != nil {
			return nil, "", err
		}
		return resp.Applications, resp.NextPageToken, nil
	}
	fetch := func(pageSize int, pageToken string) (string, error) {
		items, nextPageToken, err := it.InternalFetch(pageSize, pageToken)
		if err != nil {
			return "", err
		}
		it.items = append(it.items, items...)
		return nextPageToken, nil
	}
	it.pageInfo, it.nextFunc = iterator.NewPageInfo(fetch, it.bufLen, it.takeBuf)
	it.pageInfo.MaxSize = int(req.PageSize)
	it.pageInfo.Token = req.PageToken
	return it
}

// ApplicationIterator manages a stream of *talentpb.Application.
type ApplicationIterator struct {
	items    []*talentpb.Application
	pageInfo *iterator.PageInfo
	nextFunc func() error

	// InternalFetch is for use by the Google Cloud Libraries only.
	// It is not part of the stable interface of this package.
	//
	// InternalFetch returns results from a single call to the underlying RPC.
	// The number of results is no greater than pageSize.
	// If there are no more results, nextPageToken is empty and err is nil.
	InternalFetch func(pageSize int, pageToken string) (results []*talentpb.Application, nextPageToken string, err error)
}

// PageInfo supports pagination. See the google.golang.org/api/iterator package for details.
func (it *ApplicationIterator) PageInfo() *iterator.PageInfo {
	return it.pageInfo
}

// Next returns the next result. Its second return value is iterator.Done if there are no more
// results. Once Next returns Done, all subsequent calls will return Done.
func (it *ApplicationIterator) Next() (*talentpb.Application, error) {
	var item *talentpb.Application
	if err := it.nextFunc(); err != nil {
		return item, err
	}
	item = it.items[0]
	it.items = it.items[1:]
	return item, nil
}

func (it *ApplicationIterator) bufLen() int {
	return len(it.items)
}

func (it *ApplicationIterator) takeBuf() interface{} {
	b := it.items
	it.items = nil
	return b
}