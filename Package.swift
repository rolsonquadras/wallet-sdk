// swift-tools-version:5.5
// The swift-tools-version declares the minimum version of Swift required to build this package.

/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import PackageDescription

let version = "v0.1.5-spm"
let moduleName = "walletsdk"
let checksum = "ea2a531c3d92b93f9458d348624c677a8b6b9a1f98ff5a2c0fc4a34a38faf3c4"

let package = Package(
    name: moduleName,
    products: [
        .library(
            name: moduleName,
            targets: [moduleName]
        )
    ],
    targets: [
        .binaryTarget(
            name: moduleName,
            url: "https://github.com/rolsonquadras/wallet-sdk-test/releases/download/\(version)/\(moduleName).xcframework.zip",
            checksum: checksum
        )
    ]
)