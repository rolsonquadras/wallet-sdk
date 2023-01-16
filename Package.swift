// swift-tools-version:5.5
// The swift-tools-version declares the minimum version of Swift required to build this package.

/*
Copyright Avast Software. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import PackageDescription

let version = "0.1.2-swift-pm"
let moduleName = "walletsdk"
let checksum = "086b85180b384f389ffc6d860ad17ec1c3aaf82d61fb1a5f56b9f38c428899f5"

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