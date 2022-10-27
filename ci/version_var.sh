#
# Copyright Avast Software. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#


# Release Parameters
BASE_VERSION=0.1.0
IS_RELEASE=false

# Project Parameters
BASE_PKG_NAME=wallet-sdk
RELEASE_REPO=https://maven.pkg.github.com/trustbloc/wallet-sdk
SNAPSHOT_REPO=https://maven.pkg.github.com/trustbloc-cicd/wallet-sdk

if [ ${IS_RELEASE} = false ]
then
  EXTRA_VERSION=snapshot-$(git rev-parse --short=7 HEAD)
  PROJECT_VERSION=${BASE_VERSION}-${EXTRA_VERSION}
  PROJECT_PKG_REPO=${SNAPSHOT_REPO}
else
  PROJECT_VERSION=${BASE_VERSION}
  PROJECT_PKG_REPO=${RELEASE_REPO}
fi

export ANDROID_GROUP_ID=dev.trustbloc
export ANDROID_ARTIFACT_ID=vc-wallet-sdk
export ANDROID_VERSION=$PROJECT_VERSION