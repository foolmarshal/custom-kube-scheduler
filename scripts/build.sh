#!/usr/bin/env bash

# Copyright 2020 The OpenEBS Authors. All rights reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

#
# This script builds the application from source for multiple platforms.
set -e



# NOT NEEDED FOR NOW
#Get the parent directory of where this script is.
#SOURCE="${BASH_SOURCE[0]}"
#while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
#DIR="$( cd -P "$( dirname "$SOURCE" )/../" && pwd )"

# Change into that directory
#cd "$DIR"


GIT_COMMIT="$(git rev-parse --short HEAD)"

# Set BUILDMETA based on travis tag
#if [[ -n "$TRAVIS_TAG" ]] && [[ $TRAVIS_TAG != *"RC"* ]]; then
#    echo "released" > BUILDMETA
#fi

# Get the version details
VERSION_META="$(cat $PWD/BUILDMETA)"

# Determine the arch/os combos we're building for
UNAME=$(uname)
ARCH=$(uname -m)
if [ "$UNAME" != "Linux" -a "$UNAME" != "Darwin" ] ; then
    echo "Sorry, this OS is not supported yet."
    exit 1
fi

if [ "$UNAME" = "Darwin" ] ; then
  XC_OS="darwin"
elif [ "$UNAME" = "Linux" ] ; then
  XC_OS="linux"
fi

if [ "${ARCH}" = "i686" ] ; then
    XC_ARCH='386'
elif [ "${ARCH}" = "x86_64" ] ; then
    XC_ARCH='amd64'
elif [ "${ARCH}" = "aarch64" ] ; then
    XC_ARCH='arm64'
elif [ "${ARCH}" = "arm64" ] ; then
    XC_ARCH='arm64'
elif [ "${ARCH}" = "amd64" ] ; then
    XC_ARCH='amd64'
elif [ "${ARCH}" = "ppc64le" ] ; then
    XC_ARCH='ppc64le'
else
    echo "Unusable architecture: ${ARCH}"
    exit 1
fi


if [ -z "${PNAME}" ];
then
    echo "Project name not defined"
    exit 1
fi

if [ -z "${CTLNAME}" ];
then
    echo "CTLNAME not defined"
    exit 1
fi

# Delete the old dir
echo "==> Removing old directory..."
rm -rf bin/"${OPENEBS}"/*
mkdir -p bin/"${OPENEBS}"/

# If its dev mode, only build for ourself
if [[ "${DEV}" ]]; then
    XC_OS=$(go env GOOS)
    XC_ARCH=$(go env GOARCH)
fi

# Build!
echo "==> Building ${CTLNAME} using $(go version)... "

GOOS="${XC_OS}"
GOARCH="${XC_ARCH}"
output_name="bin/${PNAME}/"$GOOS"_"$GOARCH"/"$CTLNAME

if [ $GOOS = "windows" ]; then
    output_name+='.exe'
fi

# Build & set the Version variable in binary
# Ref: https://www.digitalocean.com/community/tutorials/using-ldflags-to-set-version-information-for-go-applications
env GOOS=$GOOS GOARCH=$GOARCH go build -ldflags \
    "-X github.com/openebs/openebsctl/cmd/openebs.Version=$GIT_COMMIT" \
    -o $output_name
#-o $optput_name $(PKG)/cmd/$


# Move all the compiled things to the $GOPATH/bin
GOPATH=${GOPATH:-$(go env GOPATH)}
case $(uname) in
    CYGWIN*)
        GOPATH="$(cygpath $GOPATH)"
        ;;
esac
OLDIFS=$IFS
IFS=: MAIN_GOPATH=($GOPATH)
IFS=$OLDIFS

# Create the gopath bin if not already available
mkdir -p ${MAIN_GOPATH}/bin/


# Copy our OS/Arch to the bin/ directory
DEV_PLATFORM="./bin/${PNAME}/$(go env GOOS)_$(go env GOARCH)"
for F in $(find ${DEV_PLATFORM} -mindepth 1 -maxdepth 2 -type f); do
    cp ${F} bin/${PNAME}/
    cp ${F} ${MAIN_GOPATH}/bin/
done

if [[ "x${DEV}" == "x" ]]; then
    # Zip and copy to the dist dir
    echo "==> Packaging..."
    for PLATFORM in $(find ./bin/${CLTNAME} -mindepth 1 -maxdepth 1 -type d); do
        OSARCH=$(basename ${PLATFORM})
        echo "--> ${OSARCH}"

        pushd "$PLATFORM" >/dev/null 2>&1
        zip ../${PNAME}-${OSARCH}.zip ./*
        popd >/dev/null 2>&1
    done
fi

# Done!
echo
echo "==> Results:"
ls -hl bin/${CLTNAME}/
