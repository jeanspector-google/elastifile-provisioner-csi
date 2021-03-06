/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"

	"github.com/golang/glog"

	"ecfs/log"
	_ "github.com/elastifile/emanage-go/src/emanage-client"
	"github.com/go-errors/errors"
)

func init() {
	_ = flag.Set("logtostderr", "true")
}

var (
	endpoint   = flag.String("endpoint", "unix://tmp/csi.sock", "CSI endpoint")
	driverName = flag.String("drivername", "csi-ecfsplugin", "Name of the driver")
	nodeId     = flag.String("nodeid", "", "Node id")
	mounter    = flag.String("volumemounter", "", "Default volume mounter (possible options are 'kernel', 'fuse')")
)

func main() {
	glog.V(log.HIGH_LEVEL_INFO).Info("Entering ECFS plugin")
	flag.Parse()

	err := updateConfigEkfs()
	if err != nil {
		err = errors.WrapPrefix(err, "Failed to update config map from within EKFS", 0)
		panic(err.Error())
	}

	driver := NewECFSDriver()
	driver.Run(*driverName, *nodeId, *endpoint, *mounter)

	glog.V(log.HIGH_LEVEL_INFO).Info("Exiting ECFS plugin")
	os.Exit(0)
}
