#!/bin/bash
#
# Copyright © 2021-2022 Dell Inc., or its subsidiaries. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#  http://www.apache.org/licenses/LICENSE-2.0

# verify-vcsi-driver method
function verify-vcsi-driver() {
  verify_k8s_versions "1.21" "1.26"
  verify_openshift_versions "4.10" "4.11"
  verify_namespace "${NS}"
 # verify_helm_values_version "${DRIVER_VERSION}"
  verify_namespace "${NS}"
 # verify_required_secrets "${RELEASE}-creds"
 # verify_optional_secrets "${RELEASE}-certs"
 # verify_optional_secrets "csirevproxy-tls-secret"
  verify_required_secrets "${RELEASE}-config"
  verify_alpha_snap_resources
  verify_sdc_installation
  verify_optional_replication_requirements
  verify_iscsi_installation
  verify_nvmetcp_installation
  verify_nvmefc_installation
  verify_snap_requirements
  verify_unity_protocol_installation
  verify_authorization_proxy_server
  verify_helm_3
}

# Check if the SDC is installed and the kernel module loaded
function verify_sdc_installation() {
  if [ ${NODE_VERIFY} -eq 0 ]; then
    return
  fi
  log step "Verifying the SDC installation"

  local SDC_MINION_NODES=$(run_command kubectl get nodes -o wide | grep -v -e master -e INTERNAL -e infra | awk ' { print $6; }')

  error=0
  missing=()
  for node in $SDC_MINION_NODES; do
    # check is the scini kernel module is loaded
    run_command ssh ${NODEUSER}@$node "/sbin/lsmod | grep scini" >/dev/null 2>&1
    rv=$?
    if [ $rv -ne 0 ]; then
      missing+=($node)
      error=1
      found_warning "SDC was not found on node: $node"
    fi
  done
  check_error error
}

function verify_unity_protocol_installation() {
if [ ${NODE_VERIFY} -eq 0 ]; then
    return
  fi

  log smart_step "Verifying sshpass installation.."
  SSHPASS=$(which sshpass)
  if [ -z "$SSHPASS" ]; then
   found_warning "sshpass is not installed. It is mandatory to have ssh pass software for multi node kubernetes setup."
  fi
  
   
  log smart_step "Verifying iSCSI installation" "$1"

  error=0
  for node in $MINION_NODES; do
    # check if the iSCSI client is installed
    echo
    echo -n "Enter the ${NODEUSER} password of ${node}: "
    read -s nodepassword
    echo
    echo "$nodepassword" > protocheckfile
    chmod 0400 protocheckfile
    unset nodepassword
    run_command sshpass -f protocheckfile ssh -o StrictHostKeyChecking=no ${NODEUSER}@"${node}" "cat /etc/iscsi/initiatorname.iscsi" > /dev/null 2>&1
    rv=$?
    if [ $rv -ne 0 ]; then
      error=1
      found_warning "iSCSI client is either not found on node: $node or not able to verify"
    fi
    run_command sshpass -f protocheckfile ssh -o StrictHostKeyChecking=no ${NODEUSER}@"${node}" "pgrep iscsid" > /dev/null 2>&1
    rv1=$?
    if [ $rv1 -ne 0 ]; then
      error=1
      found_warning "iscsid service is either not running on node: $node or not able to verify"
    fi
    rm -f protocheckfile
  done
  check_error error
}

function verify_optional_replication_requirements() {
  log step "Verifying Replication requirements"
  decho
  log arrow
  log smart_step "Verifying that Dell CSI Replication CRDs are available" "small"

  error=0
  # check for the CRDs. These are required for installation
  CRDS=("DellCSIReplicationGroups")
  for C in "${CRDS[@]}"; do
    # Verify if the CRD is present on the system
    run_command kubectl explain ${C} 2> /dev/null | grep "dell" --quiet
    if [ $? -ne 0 ]; then
      error=1
      found_warning "The CRD for ${C} is not installed. This needs to be installed if you are going to enable replication support"
    fi
  done
  check_error error
}