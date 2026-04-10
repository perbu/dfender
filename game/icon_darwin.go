//go:build darwin

package game

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

void setDockIcon(const void *rgba, int width, int height) {
    NSBitmapImageRep *rep = [[NSBitmapImageRep alloc]
        initWithBitmapDataPlanes:NULL
        pixelsWide:width
        pixelsHigh:height
        bitsPerSample:8
        samplesPerPixel:4
        hasAlpha:YES
        isPlanar:NO
        colorSpaceName:NSCalibratedRGBColorSpace
        bytesPerRow:width * 4
        bitsPerPixel:32];
    memcpy([rep bitmapData], rgba, width * height * 4);
    NSImage *img = [[NSImage alloc] initWithSize:NSMakeSize(width, height)];
    [img addRepresentation:rep];
    [[NSApplication sharedApplication] setApplicationIconImage:img];
}
*/
import "C"

import "unsafe"

func setMacOSDockIcon(pixels []byte, width, height int) {
	C.setDockIcon(unsafe.Pointer(&pixels[0]), C.int(width), C.int(height))
}
