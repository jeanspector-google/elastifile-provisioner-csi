package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/go-errors/errors"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"ecfs/log"
	"github.com/elastifile/efaasclient"
	"github.com/elastifile/efaasclient/efaasapi"
	"github.com/elastifile/efaasclient/size"
)

const (
	efaasSnapshotClassParam_Retention = "retention"
	efaasMaxSnapshotNameLen           = 63 // eFaaS API limit
	efaasTimeoutShort                 = 5 * time.Minute
	efaasTimeoutNormal                = 15 * time.Minute
	efaasTimeoutLong                  = 60 * time.Minute
)

func newEfaasClient() (efaasConf *efaasclient.Client) {
	_, secret, err := GetPluginSettings()
	if err != nil {
		panic("Failed to get plugin settings - " + err.Error())
	}

	jsonData := secret[efaasSecretsKeySaJson]
	opts := efaasclient.ClientCreateOpts{
		ProjectNumber: os.Getenv(envProjectNumber),
		BaseURL:       os.Getenv(envEfaasUrl),
	}

	if opts.ProjectNumber == "" { // TODO: DELME in v1.0.0
		opts.ProjectNumber = os.Getenv(EnvProjectNumberPrev)
	}

	efaasConf, err = efaasclient.NewClient(jsonData, opts)
	if err != nil {
		panic(fmt.Sprintf("Failed to get eFaaS client based on json %v", string(jsonData)))
	}

	return
}

func efaasGetInstanceName() string {
	return os.Getenv(envEfaasInstance)
}

// updateDefaultFsQuota works around v1 API default "auto" quota taking up all the capacity, preventing add'l filesystems' creation
func updateDefaultFsQuota(delta int64) (err error) {
	defaultFsName := efaasGetInstanceName()
	client := newEfaasClient()

	fs, err := client.GetFilesystemByName(efaasGetInstanceName(), defaultFsName)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.Warningf("Default filesystem not found - unable to update its quota")
			return nil
		}
		return errors.WrapPrefix(err, "Failed to update quota", 0)
	}

	quota := fs.HardQuota + delta
	err = client.UpdateFilesystemQuotaById(efaasGetInstanceName(), fs.Id, size.Size(quota), 5*time.Minute)
	if err != nil {
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to update default filesystem quota to %v", quota), 0)
	}

	return
}

func efaasCreateEmptyVolume(volOptions *volumeOptions) (volumeId volumeHandleType, err error) {
	client := newEfaasClient()
	glog.V(log.DETAILED_INFO).Infof("ecfs: Creating Volume - settings: %+v", volOptions)
	volumeId = volOptions.VolumeId

	snapshot := &efaasapi.SnapshotSchedule{
		Enable:    false,
		Schedule:  "Monthly",
		Retention: 2.0,
	}

	var accessorItems []efaasapi.AccessorItems

	err = json.Unmarshal([]byte(volOptions.ClientRules), &accessorItems) // Has to be JSON due to AccessorItems only having json tags
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to unmarshal client rules JSON: %v", volOptions.ClientRules), 0)
		return
	}

	filesystem := efaasapi.DataContainerAdd{
		Name:        string(volumeId),
		HardQuota:   volOptions.Capacity,
		QuotaType:   efaasclient.QuotaTypeFixed,
		Description: fmt.Sprintf("Filesystem %v", volumeId),
		Accessors:   &efaasapi.Accessors{Items: accessorItems},
		Snapshot:    snapshot,
	}

	// Create Filesystem
	err = client.AddFilesystem(efaasGetInstanceName(), filesystem, efaasTimeoutNormal)
	if err != nil {
		if isErrorAlreadyExists(err) {
			glog.V(log.DEBUG).Infof("ecfs: Volume %v was already created - assuming it was created "+
				"during previous, failed, attempt", volumeId)
			_, e := client.GetFilesystemByName(efaasGetInstanceName(), string(volumeId))
			if e != nil {
				logSecondaryError(err, e)
				return
			}
			err = nil // This error is acceptable
		} else {
			quotaError := "Quota exceeded"
			if strings.Contains(err.Error(), quotaError) {
				// Resize default filesystem
				e := updateDefaultFsQuota(-1 * volOptions.Capacity)
				if e != nil {
					logSecondaryError(err, e)
					return
				}

				// Retry AddFilesystem
				e = client.AddFilesystem(efaasGetInstanceName(), filesystem, efaasTimeoutNormal)
				if e != nil {
					logSecondaryError(err, e)
					return
				}

				err = nil // This error has been worked around
			} else {
				err = errors.Wrap(err, 0)
				return "", errors.Wrap(err, 0)
			}
		}
	}
	//volOptions.DataContainer = fs

	glog.V(log.DEBUG).Infof("ecfs: Created volume with id %v", volumeId)

	return
}

func efaasDeleteVolume(volName volumeHandleType) (err error) {
	client := newEfaasClient()
	err = client.DeleteFilesystem(efaasGetInstanceName(), string(volName), efaasTimeoutLong)
	if err != nil {
		if isErrorDoesNotExist(err) {
			glog.V(log.DEBUG).Infof("ecfs: Filesystem %v not found - assuming already deleted", volName)
			return nil
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to delete filesystem %v", volName), 0)
	}

	glog.V(log.DETAILED_INFO).Infof("ecfs: Deleted filesystem %v", volName)
	return nil
}

const (
	efaasSnapshotStatus_PENDING = "PENDING"
	efaasSnapshotStatus_READY   = "READY"
)

func efaasIsSnapshotUsable(snapshot *efaasapi.Snapshots) (isUsable bool) {
	isUsable = snapshot.Status == efaasSnapshotStatus_READY
	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Is snapshot %v usable? %v - %#v ", snapshot.Name, isUsable, snapshot)
	return
}

// efaasParseTimestamp parses dateTime (e.g. 2019-05-27 13:03:12) into protobuf timestamp
func efaasParseTimestamp(dateTime string) (ts *timestamp.Timestamp, err error) {
	return parseTimestamp(dateTime, "2006-01-02 15:04:05")
}

func efaasGetCreateSnapshotResponse(efaasSnapshot *efaasapi.Snapshots, req *csi.CreateSnapshotRequest) (response *csi.CreateSnapshotResponse, err error) {
	creationTimestamp, err := efaasParseTimestamp(efaasSnapshot.CreationTimestamp)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	isReady := efaasIsSnapshotUsable(efaasSnapshot)

	response = &csi.CreateSnapshotResponse{
		Snapshot: &csi.Snapshot{
			SnapshotId:     efaasSnapshot.Name,
			SourceVolumeId: req.GetSourceVolumeId(),
			CreationTime:   creationTimestamp,
			ReadyToUse:     isReady,
		},
	}

	return
}

func efaasGetSnapshotById(snapshotId string) (snapshot efaasapi.Snapshots, err error) {
	client := newEfaasClient()
	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Getting Snapshot by ID: %v", snapshotId)
	snapshot, err = client.GetSnapshotById(snapshotId)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get eFaaS snapshot by ID %v", snapshotId), 0)
		return
	}
	return
}

func efaasGetSnapshotByName(snapshotName string) (snapshot efaasapi.Snapshots, err error) {
	client := newEfaasClient()
	glog.V(log.VERBOSE_DEBUG).Infof("ecfs: Getting Snapshot by name: %v", snapshotName)
	snapshot, err = client.GetSnapshotByName(efaasGetInstanceName(), snapshotName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get eFaaS snapshot by name %v", snapshotName), 0)
		return
	}
	return
}

func efaasSnapshotRetentionFromParams(params map[string]string) (retention float32, err error) {
	retStr, found := params[efaasSnapshotClassParam_Retention]
	if !found {
		err = errors.Errorf("Parameter %v not found on volume snapshot class parameters: %v",
			efaasSnapshotClassParam_Retention, params)
		return
	}

	ret64, err := strconv.ParseFloat(retStr, 32)
	if !found {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to parse string %v to float", retStr), 0)
		return
	}
	retention = float32(ret64)
	return float32(ret64), err
}

func efaasCreateSnapshot(name string, volumeId volumeHandleType, params map[string]string) (snapshot efaasapi.Snapshots, err error) {
	client := newEfaasClient()
	glog.V(log.HIGH_LEVEL_INFO).Infof("ecfs: Creating snapshot %v for volume %v", name, volumeId)
	glog.V(log.DEBUG).Infof("ecfs: Creating snapshot %v - parameters: %v", name, params)

	fsName := string(volumeId)
	retention, err := efaasSnapshotRetentionFromParams(params)
	if err != nil {
		err = errors.Wrap(err, 0)
		return
	}

	if len(name) > efaasMaxSnapshotNameLen {
		err = errors.Errorf("Requested snapshot name length (%v = %v) is over the max supported length (%v)",
			name, len(name), efaasMaxSnapshotNameLen)
		return
	}

	snapCreateArgs := efaasapi.Snapshot{
		Name:      name,
		Retention: retention,
	}

	err = client.CreateSnapshot(efaasGetInstanceName(), fsName, snapCreateArgs, efaasTimeoutShort)
	if err != nil {
		if !isErrorAlreadyExists(err) {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create snapshot on filesystem %v - %#v",
				volumeId, snapCreateArgs), 0)
			return
		}
	}

	snapshot, err = client.GetSnapshotByFsAndName(efaasGetInstanceName(), fsName, name)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get snapshot by name %v from filesystem %v",
			name, fsName), 0)
		return
	}

	return
}

// efaasCreateShare creates share (aka export) on the specified snapshot
func efaasCreateShare(snapName string) (share *efaasapi.Share, err error) {
	client := newEfaasClient()
	err = client.CreateShare(efaasGetInstanceName(), snapName, snapshotExportName, efaasTimeoutShort)
	if err != nil {
		if !isErrorAlreadyExists(err) {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create share on snapshot %v", snapName), 0)
			return
		}
	}

	share, err = client.GetShare(efaasGetInstanceName(), snapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to get recently created share %v on snapshot %v",
			snapshotExportName, snapName), 0)
		return
	}

	return
}

func efaasDeleteShare(snapName string) (err error) {
	client := newEfaasClient()
	err = client.DeleteShare(efaasGetInstanceName(), snapName, snapshotExportName, efaasTimeoutShort)
	if err != nil {
		if !isErrorDoesNotExist(err) {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete share %v on snapshot %v",
				snapshotExportName, snapName), 0)
			return
		}
		return nil
	}

	return
}

// efaasCreateVolumeFromSnapshot is intended to be used by snapshot restore/clone functions
func efaasCreateVolumeFromSnapshot(srcSnapName string, dstVolOptions *volumeOptions) (dstVolumeId volumeHandleType, err error) {
	var srcSnapMountPath = fmt.Sprintf("/mnt/%v", srcSnapName)

	glog.V(log.DETAILED_INFO).Infof("ecfs: Restoring snapshot %v - dstVolOptions: %+v", srcSnapName, dstVolOptions)

	srcSnapshot, err := efaasGetSnapshotByName(srcSnapName)
	if err != nil {
		err = status.Error(codes.Internal, errors.WrapPrefix(err,
			fmt.Sprintf("Failed to to get source snapshot by name %v", srcSnapName), 0).Error())
		return
	}

	// Create snapshot export
	_, err = efaasCreateShare(srcSnapName)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create share on snapshot %v", srcSnapName), 0)
		return
	}

	// Create destination volume
	dstVolumeId, err = efaasCreateEmptyVolume(dstVolOptions)
	if err != nil {
		if !isErrorAlreadyExists(err) {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to create destination volume %v",
				dstVolOptions.VolumeId), 0)
			glog.Errorf(err.Error())
			err = status.Error(codes.Internal, err.Error())
			return
		} else {
			glog.V(log.DEBUG).Infof("ecfs: Destination volume %v exists. Assuming CSI retried operation "+
				"that was previous aborted", dstVolOptions.VolumeId)
		}
	}

	// Mount the source snapshot
	err = mountEfaasSnapshot(srcSnapMountPath, srcSnapshot.Name)
	if err != nil {
		isMount, e := isMountPoint(srcSnapMountPath) // TODO: Consider remounting or using unique paths
		if e != nil {
			logSecondaryError(err, e)
		}
		if !isMount {
			err = errors.WrapPrefix(err, "Failed to mount source snapshot's export", 0)
			return
		} else {
			glog.V(log.DEBUG).Infof("ecfs: Source snapshot is already mounted on %v - "+
				"assuming it was mounted during previous, aborted, attempt", srcSnapMountPath)
		}
	}

	defer func() { // Umount the source export
		e := unmountAndCleanup(srcSnapMountPath)
		if e != nil {
			if err == nil {
				err = errors.WrapPrefix(e, "Failed to unmount source snapshot", 0)
				glog.Warning(err.Error())
			} else {
				logSecondaryError(err, e)
			}
		}
	}()

	// Mount the destination volume
	dstVolMountPath := fmt.Sprintf("/mnt/%v", dstVolumeId)
	err = mountEcfs(dstVolMountPath, dstVolumeId, []string{"vers=3"})
	if err != nil {
		isMount, e := isMountPoint(dstVolMountPath) // TODO: Consider remounting or using unique paths
		if e != nil {
			logSecondaryError(err, e)
		}
		if !isMount {
			err = errors.WrapPrefix(err, "Failed to mount destination volume", 0)
			return
		} else {
			glog.V(log.DEBUG).Infof("ecfs: Destination volume is already mounted on %v - "+
				"assuming it was mounted during previous, aborted, attempt", dstVolMountPath)
		}
	}

	defer func() { // Umount the destination volume
		e := unmountAndCleanup(dstVolMountPath)
		if e != nil {
			if err == nil {
				err = errors.WrapPrefix(e, "Failed to unmount destination volume", 0)
				glog.Warning(err.Error())
			} else {
				logSecondaryError(err, e)
			}
		}
	}()

	// Copy the source snapshot's contents into the destination volume
	err = copyDirWithKeepalive(srcSnapMountPath, dstVolMountPath, string(dstVolumeId))
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to copy snapshot %v (%v) contents to volume %v (%v)",
			srcSnapName, srcSnapMountPath, dstVolumeId, dstVolMountPath), 0)
		return
	}

	// Testing instrumentation
	delaySec := getDebugValueInt(debugValueCloneDelaySec, nil)
	if delaySec > 0 {
		glog.V(log.VERBOSE_DEBUG).Infof("ecfs: DEBUG - delaying volume creation completion by %v sec", delaySec)
		time.Sleep(time.Duration(delaySec) * time.Second)
	}

	return
}

func efaasDeleteSnapshot(name string) (err error) {
	client := newEfaasClient()

	glog.V(log.INFO).Infof("ecfs: Deleting snapshot %v", name)

	snapshot, err := client.GetSnapshotByName(efaasGetInstanceName(), name)
	if err != nil {
		if isErrorDoesNotExist(err) { // This operation has to be idempotent
			glog.V(log.DEBUG).Infof("ecfs: Snapshot %v not found - assuming already deleted", name)
			return nil
		}
	}

	// Delete share
	if snapshot.Share.Name != "" {
		err = efaasDeleteShare(name)
		if err != nil {
			err = errors.WrapPrefix(err, fmt.Sprintf("Failed to delete share from snapshot %v", name), 0)
			return
		}
	}

	// Delete snapshot
	err = client.DeleteSnapshot(efaasGetInstanceName(), snapshot.FilesystemName, name, efaasTimeoutShort)
	if err != nil {
		if isErrorDoesNotExist(err) { // This operation has to be idempotent
			glog.V(log.DEBUG).Infof("ecfs: Snapshot %v not found - assuming already deleted", name)
			return nil
		}
		if isWorkaround("EL-13618 - Failed read-dir") {
			const EL13618 = "Failed read-dir"
			if strings.Contains(err.Error(), EL13618) {
				glog.Warningf("ecfs: Snapshot delete failed due to EL-13618 - returning success to cleanup the pv. Actual error: %v", err)
				return nil
			}
		}
		return errors.WrapPrefix(err, fmt.Sprintf("Failed to delete snapshot %v", name), 0)
	}

	return
}

func efaasCloneVolume(source *csi.VolumeContentSource_VolumeSource, dstVolOptions *volumeOptions) (dstVolumeId volumeHandleType, err error) {
	var (
		reqParams   map[string]string
		srcVolumeId = volumeHandleType(source.GetVolumeId())
		srcSnapName = truncateStr(fmt.Sprintf("4-%v", dstVolOptions.VolumeId), maxSnapshotNameLen)
	)

	glog.V(log.DETAILED_INFO).Infof("ecfs: Cloning volume %v to %v via snapshot %v - dstVolOptions: %+v",
		srcVolumeId, dstVolOptions.VolumeId, srcSnapName, dstVolOptions)

	// Take source volume's snapshot
	_, err = efaasCreateSnapshot(srcSnapName, srcVolumeId, reqParams)
	if err != nil {
		err = errors.WrapPrefix(err,
			fmt.Sprintf("Failed to create snapshot for volume %v with name %v", srcVolumeId, srcSnapName), 0)
		return
	}

	defer func() { // Cleanup source snapshot
		e := efaasDeleteSnapshot(srcSnapName)
		if e != nil {
			if err == nil {
				err = errors.WrapPrefix(e, fmt.Sprintf("Failed to delete source snapshot %v", srcSnapName), 0)
				glog.Warning(e.Error())
			} else {
				logSecondaryError(err, e)
			}
		}
	}()

	dstVolumeId, err = efaasCreateVolumeFromSnapshot(srcSnapName, dstVolOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to clone volume %v", srcVolumeId), 0)
		return
	}

	return
}

func efaasRestoreSnapshotToVolume(source *csi.VolumeContentSource_SnapshotSource, dstVolOptions *volumeOptions) (dstVolumeId volumeHandleType, err error) {
	dstVolumeId, err = efaasCreateVolumeFromSnapshot(source.GetSnapshotId(), dstVolOptions)
	if err != nil {
		err = errors.WrapPrefix(err, fmt.Sprintf("Failed to restore snapshot %v", source.GetSnapshotId()), 0)
		return
	}

	return
}
