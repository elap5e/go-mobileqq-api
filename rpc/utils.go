package rpc

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"

	"github.com/elap5e/go-mobileqq-api/tlv"
)

var (
	deviceGUID          = md5.Sum(append(defaultDeviceOSBuildID, defaultDeviceMACAddress...)) // []byte("%4;7t>;28<fclient.5*6")
	deviceGUIDFlag      = uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00))
	deviceIsGUIDFileNil = false
	deviceIsGUIDGenSucc = true
	deviceIsGUIDChanged = false
	deviceDPWD          = []byte{}

	clientVerifyMethod = uint8(0x82) // 0x00, 0x82
)

var (
	clientPackageName  = []byte("com.tencent.mobileqq")
	clientVersionName  = []byte("8.8.3")
	clientRevision     = "8.8.3.b2791edc"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID      = uint32(0x20030cb2)
	clientBuildTime  = uint64(0x00000000609b85ad)
	clientSDKVersion = "6.0.0.2476"
	clientSSOVersion = uint32(0x00000011)

	clientCodecAppIDDebug   = []byte("736350642")
	clientCodecAppIDRelease = []byte("736350642")
)

var (
	clientImageType  = uint8(0x01)
	clientMiscBitmap = uint32(0x08f7ff7c)
)

func init() {
	deviceDPWD = func(n int) []byte {
		b := make([]byte, n)
		for i := range b {
			b[i] = byte(0x41 + rand.Intn(1)*0x20 + rand.Intn(26))
		}
		return b
	}(16)
	log.Printf("--> [init] dump device dpwd\n%s", hex.Dump(deviceDPWD))
}

func SetClientForAndroidPhone() {
	clientPackageName = []byte("com.tencent.mobileqq")
	clientVersionName = []byte("8.8.3")
	clientRevision = "8.8.3.b2791edc"
	clientSignatureMD5 = [16]byte{0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77, 0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d}

	clientAppID = uint32(0x20030cb2)
	clientBuildTime = uint64(0x00000000609b85ad)
	clientSDKVersion = "6.0.0.2476"
	clientSSOVersion = uint32(0x00000011)

	clientCodecAppIDDebug = []byte("736350642")
	clientCodecAppIDRelease = []byte("736350642")
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientForAndroidPad() {
	clientPackageName = []byte("com.tencent.minihd.qq")
	clientVersionName = []byte("5.9.2")
	clientRevision = "5.9.2.3baec0"
	clientSignatureMD5 = [16]byte{0xaa, 0x39, 0x78, 0xf4, 0x1f, 0xd9, 0x6f, 0xf9, 0x91, 0x4a, 0x66, 0x9e, 0x18, 0x64, 0x74, 0xc7}

	clientAppID = uint32(0x2002fdd5)
	clientBuildTime = uint64(0x000000005f1e8730)
	clientSDKVersion = "6.0.0.2433"
	clientSSOVersion = uint32(0x0000000c)

	clientCodecAppIDDebug = []byte("73636270;")
	clientCodecAppIDRelease = []byte("736346857")
	tlv.SetSSOVersion(clientSSOVersion)
}

func SetClientForiPhone() {
	panic("not implement")
}

func SetClientForiPad() {
	panic("not implement")
}
