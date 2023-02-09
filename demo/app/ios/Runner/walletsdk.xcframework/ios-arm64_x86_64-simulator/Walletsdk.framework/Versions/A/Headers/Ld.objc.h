// Objective-C API for talking to github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/ld Go package.
//   gobind -lang=objc github.com/trustbloc/wallet-sdk/cmd/wallet-sdk-gomobile/ld
//
// File is generated by gobind. Do not edit.

#ifndef __Ld_H__
#define __Ld_H__

@import Foundation;
#include "ref.h"
#include "Universe.objc.h"

#include "Walleterror.objc.h"
#include "Api.objc.h"

@class LdDocLoader;

/**
 * DocLoader represents a type that can help with linked domains.
 */
@interface LdDocLoader : NSObject <goSeqRefInterface, ApiLDDocumentLoader> {
}
@property(strong, readonly) _Nonnull id _ref;

- (nonnull instancetype)initWithRef:(_Nonnull id)ref;
/**
 * NewDocLoader returns a new DocLoader instance.
 */
- (nullable instancetype)init;
/**
 * LoadDocument load linked document by url.
 */
- (ApiLDDocument* _Nullable)loadDocument:(NSString* _Nullable)u error:(NSError* _Nullable* _Nullable)error;
@end

/**
 * NewDocLoader returns a new DocLoader instance.
 */
FOUNDATION_EXPORT LdDocLoader* _Nullable LdNewDocLoader(void);

#endif
