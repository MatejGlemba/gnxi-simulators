# SPDX-FileCopyrightText: 2022 2020-present Open Networking Foundation <info@opennetworking.org>
#
# SPDX-License-Identifier: Apache-2.0

FROM onosproject/golang-build:v1.2.0 as build

RUN cd $GOPATH && GO111MODULE=on go install github.com/google/gnxi/gnoi_target@c86e52276d175750acd8e5abc62dd35f5393713f
RUN cd $GOPATH && GO111MODULE=on go install github.com/google/gnxi/gnoi_cert@c86e52276d175750acd8e5abc62dd35f5393713f

ENV ONOS_SIMULATORS_ROOT=/go/src/github.com/onosproject/gnxi-simulators
ENV CGO_ENABLED=0

RUN mkdir -p $ONOS_SIMULATORS_ROOT/

COPY . $ONOS_SIMULATORS_ROOT

RUN cd $ONOS_SIMULATORS_ROOT && GO111MODULE=on go build -o /go/bin/gnmi_target ./cmd/gnmi_target


FROM alpine:3.11
RUN apk add bash openssl curl libc6-compat
ENV GNMI_PORT=10161
ENV GNMI_INSECURE_PORT=11161
ENV GNOI_PORT=50001
ENV SIM_MODE=1
ENV HOME=/home/devicesim
RUN mkdir $HOME
WORKDIR $HOME

COPY --from=build /go/bin/gn* /usr/local/bin/

COPY configs/target_configs target_configs
COPY tools/scripts scripts
COPY pkg/certs certs

RUN chmod +x ./scripts/run_target.sh
CMD ["./scripts/run_target.sh"]
