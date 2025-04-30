# Copyright 2025 iLogtail Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.24-alpine AS build

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

RUN apk add --no-cache git

WORKDIR /src
COPY . .

RUN go build -o ConfigServer -ldflags="-s -w"

RUN chmod 755 /src/ConfigServer && \
    sed -i 's/127\.0\.0\.1/0\.0\.0\.0/' /src/setting/setting.json

FROM scratch

ENV GIN_MODE=release

WORKDIR /config_server
COPY --from=build /src/ConfigServer /config_server/ConfigServer
COPY --from=build /src/setting/setting.json /config_server/setting/setting.json

CMD ["/config_server/ConfigServer"]
