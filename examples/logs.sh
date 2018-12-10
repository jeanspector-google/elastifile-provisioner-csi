#!/usr/bin/env bash

: ${POD_ID:=2}
: ${CONTAINER_NAME:=csi-ecfsplugin}
POD_NAME=$(kubectl get pods -l app=$CONTAINER_NAME -o=name | head -n ${POD_ID} | tail -n 1)

function get_pod_status() {
	echo -n $(kubectl get $POD_NAME -o jsonpath="{.status.phase}")
}

while [[ "$(get_pod_status)" != "Running" ]]; do
	sleep 1
	echo "Waiting for $POD_NAME (status $(get_pod_status))"
done

kubectl logs -f $POD_NAME -c $CONTAINER_NAME
