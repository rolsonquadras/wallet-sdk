!/bin/bash
#
# Copyright Gen Digital Inc. All Rights Reserved.
#
# SPDX-License-Identifier: Apache-2.0
#

echo "Running $0"

cd demo/app
pwd

buildNumber=5

# Dev Setup
bundleIdentifier=com.securekey.trustbloc.vcwallet.dev
bundleDisplayName=TrustBloc-Dev
iconFileName=app-icon--dev.jpg

# STG Setup
#bundleIdentifier=com.securekey.trustbloc.vcwallet.stg
#bundleDisplayName=TrustBloc-STG
#iconFileName=app-icon--stg.jpg

#sed -i -E "s/appicon.png/app-icon--dev.jpg/g" pubspec.yaml
sed -i "s/1.0.0+1/1.0.0+$buildNumber/g" pubspec.yaml
sed -i "s/appicon.png/$iconFileName/g" pubspec.yaml

flutter clean
flutter pub get

flutter pub run flutter_launcher_icons:main

cd ios/
pod install

sed -i 's/readlink /readlink -f /g' Pods/Target\ Support\ Files/Pods-Runner/Pods-Runner-frameworks.sh

sed -i "s/TrustBloc/$bundleDisplayName/g" Runner/Info.plist
sed -i "s/\$(FLUTTER_BUILD_NUMBER)/$buildNumber/g" Runner/Info.plist
sed -i "s/\$(PRODUCT_BUNDLE_IDENTIFIER)/$bundleIdentifier/g" Runner/Info.plist

sed -i "s/TrustBloc/$bundleDisplayName/g" Runner.xcodeproj/project.pbxproj
sed -i "s/\"\$(FLUTTER_BUILD_NUMBER)\"/5/g" Runner.xcodeproj/project.pbxproj
sed -i "s/dev.trustbloc.wallet/$bundleIdentifier/g" Runner.xcodeproj/project.pbxproj
