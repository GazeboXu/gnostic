@REM #!/bin/sh
@REM #
@REM # Copyright 2016 Google LLC. All Rights Reserved.
@REM #
@REM # Licensed under the Apache License, Version 2.0 (the "License");
@REM # you may not use this file except in compliance with the License.
@REM # You may obtain a copy of the License at
@REM #
@REM #    http://www.apache.org/licenses/LICENSE-2.0
@REM #
@REM # Unless required by applicable law or agreed to in writing, software
@REM # distributed under the License is distributed on an "AS IS" BASIS,
@REM # WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
@REM # See the License for the specific language governing permissions and
@REM # limitations under the License.
@REM #

go get github.com/golang/protobuf/protoc-gen-go

@REM protoc --go_out=. openapiv2/OpenAPIv2.proto
@REM protoc --go_out=. openapiv3/OpenAPIv3.proto
@REM protoc --go_out=. discovery/discovery.proto
@REM protoc --go_out=. plugins/plugin.proto
@REM protoc --go_out=. extensions/extension.proto
@REM protoc --go_out=. surface/surface.proto
@REM protoc --go_out=. metrics/vocabulary.proto
@REM protoc --go_out=. metrics/complexity.proto
