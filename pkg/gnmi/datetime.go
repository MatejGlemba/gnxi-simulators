// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gnmi

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/golang/protobuf/proto"
	pb "github.com/openconfig/gnmi/proto/gnmi"
)

// SetDateTime update current-datetime field in runtime
func (s *Server) SetDateTime() error {
	s.configMu.Lock()
	defer s.configMu.Unlock()
	var path pb.Path
	textPbPath := `elem:<name:"system" > elem:<name:"state" > elem:<name:"current-datetime" > `
	if err := proto.UnmarshalText(textPbPath, &path); err != nil {
		return err
	}

	val := &pb.TypedValue{
		Value: &pb.TypedValue_StringVal{
			StringVal: time.Now().Format("2006-01-02T15:04:05Z-07:00"),
		},
	}
	update := &pb.Update{Path: &path, Val: val}

	jsonTree, _ := ygot.ConstructIETFJSON(s.config, &ygot.RFC7951JSONConfig{})
	_, _ = s.doReplaceOrUpdate(jsonTree, pb.UpdateResult_UPDATE, nil, update.GetPath(), update.GetVal())
	jsonDump, err := json.Marshal(jsonTree)
	if err != nil {
		msg := fmt.Sprintf("error in marshaling IETF JSON tree to bytes: %v", err)
		log.Error(msg)
		return status.Error(codes.Internal, msg)
	}
	rootStruct, err := s.model.NewConfigStruct(jsonDump)
	if err != nil {
		msg := fmt.Sprintf("error in creating config struct from IETF JSON data: %v", err)
		log.Error(msg)
		return status.Error(codes.Internal, msg)
	}
	s.config = rootStruct
	s.ConfigUpdate <- update
	return nil

}
