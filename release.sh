#!/usr/bin/env bash
#
# Copyright 2021 OpenIndex.de.
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
#

NAME="JavaMacLauncher"
MAC_VERSION="10.9"

export CGO_CFLAGS="-mmacosx-version-min=${MAC_VERSION}"
export CGO_LDFLAGS="-mmacosx-version-min=${MAC_VERSION} -s -w"

DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC="${DIR}/src"
TARGET="${DIR}/target"

rm -Rf "${TARGET}"
mkdir -p "${TARGET}"

cd "${SRC}" || exit
"${DIR}/go.sh" build -o "${TARGET}/${NAME}" -a -v
