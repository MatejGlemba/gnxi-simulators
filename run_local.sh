#!/bin/bash

# SPDX-FileCopyrightText: 2022 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

hostname=localhost
GNMI_INSECURE_PORT=10161
GNMI_PORT=11161
HOME=/home/matej/Frinx/gnxi-simulators

if [ "$hostname" != "localhost" ]; then \
    IPADDR=`ip route get 1.2.3.4 | grep dev | awk '{print $7}'`
    $HOME/pkg/certs/generate_certs.sh $hostname > /dev/null 2>&1;
    echo "Please add '"$IPADDR" "$hostname"' to /etc/hosts and access with gNMI client at "$hostname":"$GNMI_PORT; \
else \
    echo "gNMI target in secure mode is on $hostname:"${GNMI_PORT};
    echo "gNMI target insecure mode is on $hostname:"${GNMI_INSECURE_PORT};
fi
sed -i -e "s/replace-device-name/"$hostname"/g" $HOME/configs/target_configs/typical_ofsw_config.json && \
sed -i -e "s/replace-motd-banner/Welcome to gNMI service on "$hostname":"$GNMI_PORT"/g" $HOME/configs/target_configs/typical_ofsw_config.json

build/bin/gnmi_target \
    -bind_address :$GNMI_INSECURE_PORT \
    -alsologtostderr \
    -notls \
    -insecure \
    -config $HOME/configs/target_configs/typical_ofsw_config.json &

build/bin/gnmi_target \
    -bind_address :$GNMI_PORT \
    -key $HOME/pkg/certs/$hostname.key \
    -cert $HOME/pkg/certs/$hostname.crt \
    -ca $HOME/pkg/certs/onfca.crt \
    -alsologtostderr \
    -config $HOME/configs/target_configs/typical_ofsw_config.json
