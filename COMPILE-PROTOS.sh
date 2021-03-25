#!/bin/sh
#
# Copyright 2016 Google LLC. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

go get github.com/golang/protobuf/protoc-gen-go

# 最新版本(v1.5.1)的protoc-gen-go使用了v1.26.0版的google.golang.org/protobuf，
# 会导致如下命令出错，这些也不会用到，去掉

# protoc --go_out=. openapiv2/OpenAPIv2.proto
# protoc --go_out=. openapiv3/OpenAPIv3.proto
# protoc --go_out=. discovery/discovery.proto
# protoc --go_out=. plugins/plugin.proto
# protoc --go_out=. extensions/extension.proto
# protoc --go_out=. surface/surface.proto
# protoc --go_out=. metrics/vocabulary.proto
# protoc --go_out=. metrics/complexity.proto
