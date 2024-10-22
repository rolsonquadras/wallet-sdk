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

cd android/

sed -i "s/dev.trustbloc.wallet/$bundleIdentifier/g" app/build.gradle
sed -i "s/signingConfig signingConfigs.debug/signingConfig signingConfigs.release/g" app/build.gradle

sed -i "8i def keystoreProperties = new Properties() \n \
        def keystorePropertiesFile = rootProject.file('key.properties') \n \
        if (keystorePropertiesFile.exists()) { \n \
               keystoreProperties.load(new FileInputStream(keystorePropertiesFile)) \n \
        }" app/build.gradle


sed -i "65i signingConfigs {\n \
                    release {\n \
                       keyAlias keystoreProperties['keyAlias']\n \
                       keyPassword keystoreProperties['keyPassword']\n \
                       storeFile keystoreProperties['storeFile'] ? file(keystoreProperties['storeFile']) : null\n \
                       storePassword keystoreProperties['storePassword']\n \
                    }\n \
                }" app/build.gradle

sed -i "s/TrustBloc/$bundleDisplayName/g" app/src/main/AndroidManifest.xml

flutter build appbundle