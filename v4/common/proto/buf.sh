#!/usr/bin/env bash

#
# Copyright (c) 2019-2021. Abstrium SAS <team (at) pydio.com>
# This file is part of Pydio Cells.
#
# Pydio Cells is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Pydio Cells is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Pydio Cells.  If not, see <http://www.gnu.org/licenses/>.
#
# The latest code can be found at <https://pydio.com>.
#

#cd $GOPATH/src

#protoc -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
#       -I$GOPATH/src/github.com/google/protobuf/src \
#       -I$GOPATH/src \
#       -I. \
#       --go_out=plugins=grpc:. \
#       --govalidators_out=. \
#       $GOPATH/src/github.com/pydio/cells/v4/common/proto/$1/*.proto

if [ ! -d "./bin" ]
then
  mkdir -p ./bin
  export GOBIN=$PWD/bin
  export PATH=$GOBIN:$PATH
  go get google.golang.org/protobuf/cmd/protoc-gen-go
  go get github.com/micro/micro/v3/cmd/protoc-gen-micro
  go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
  go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
  go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
  go get github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
else
  export GOBIN=$PWD/bin
  export PATH=$GOBIN:$PATH
fi

cd $1

buf generate --output $GOPATH/src

cd -