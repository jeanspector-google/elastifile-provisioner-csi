package main

/*
This file handles the config map storing values used in testing instrumentation
*/

import (
	"fmt"
	"strconv"

	"github.com/go-errors/errors"
	"github.com/golang/glog"

	"ecfs/co"
	"ecfs/log"
)

const (
	debugConfigMapName      = "csi-debug"
	debugValueCloneDelaySec = "cloneDelaySec"
	debugValueCopyDelaySec  = "copyDelaySec"
)

func logDebugValue(key string, value interface{}) {
	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Using debug value %v=%v", key, value)
}

// getDebugValue returns string value corresponding to the key in debugConfigMapName, and the default value otherwise
// If nil is specified, zero value is returned
func getDebugValue(key string, defaultValue *string) string {
	value, err := co.GetConfigMapValue(Namespace(), debugConfigMapName, key)
	if err != nil {
		if defaultValue != nil {
			value = *defaultValue
		}

		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get debug value %v from config map %v - using default",
			key, debugConfigMapName), 0)
		glog.V(log.VERBOSE_DEBUG).Infof(err.Error())
	}

	logDebugValue(key, value)
	return value
}

// getDebugValueInt returns integer value corresponding to the key in debugConfigMapName, and the default value otherwise
// If nil is specified, zero value is returned
func getDebugValueInt(key string, defaultValue *int) (value int) {
	valueStr := getDebugValue(key, nil)
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		if defaultValue != nil {
			value = *defaultValue
		}

		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to convert debug value %v from %v to int - using default",
			key, valueStr), 0)
		glog.V(log.VERBOSE_DEBUG).Infof(err.Error())
	}

	logDebugValue(key, value)
	return value
}
