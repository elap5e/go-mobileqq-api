package rpc

import (
	"crypto/md5"

	"github.com/elap5e/go-mobileqq-api/util"
)

type Config struct {
	BaseDir  string
	CacheDir string
	Client   *ClientConfig
	Device   *DeviceConfig
}

type ClientConfig struct {
	AppID             uint32
	CodecAppIDDebug   string
	CodecAppIDRelease string

	PackageName  string
	VersionName  string
	Revision     string
	SignatureMD5 []byte

	BuildTime  uint64
	SDKVersion string
	SSOVersion uint32

	ImageType  uint8
	MiscBitmap uint32

	CanCaptcha bool
}

type DeviceConfig struct {
	Bootloader   string
	ProcVersion  string
	Codename     string
	Incremental  string
	Fingerprint  string
	BootID       string
	OSBuildID    string
	Baseband     string
	InnerVersion string

	GUID          []byte
	GUIDFlag      uint32
	IsGUIDFileNil bool
	IsGUIDGenSucc bool
	IsGUIDChanged bool
}

func NewDeviceConfig() *DeviceConfig {
	return NewDeviceConfigForAndroid()
}

func NewDeviceConfigForAndroid() *DeviceConfig {
	return &DeviceConfig{
		Bootloader:   "unknown",
		ProcVersion:  "Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)",
		Codename:     "davinci",
		Incremental:  "20.10.20",
		Fingerprint:  "Xiaomi/davinci/davinci:11/RKQ1.200827.002/20.10.20:user/release-keys",
		BootID:       "aa6bf49c-a995-4761-874f-8b1a9eee341e",
		OSBuildID:    "RKQ1.200827.002",
		Baseband:     "4.3.c5-00069-SM6150_GEN_PACK-1",
		InnerVersion: "20.10.20",

		GUID:          util.STBytesTobytes(md5.Sum(append(defaultDeviceOSBuildID, defaultDeviceMACAddress...))), // []byte("%4;7t>;28<fclient.5*6")
		GUIDFlag:      uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00)),
		IsGUIDFileNil: false,
		IsGUIDGenSucc: true,
		IsGUIDChanged: false,
	}
}

func NewClientConfig() *ClientConfig {
	return NewClientConfigForAndroid()
}

func NewClientConfigForAndroid() *ClientConfig {
	return &ClientConfig{
		AppID:             0x20030cb2,
		CodecAppIDDebug:   "736350642",
		CodecAppIDRelease: "736350642",

		PackageName: "com.tencent.mobileqq",
		VersionName: "8.8.3",
		Revision:    "8.8.3.b2791edc",
		SignatureMD5: []byte{
			0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77,
			0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d,
		},

		BuildTime:  0x00000000609b85ad,
		SDKVersion: "6.0.0.2476",
		SSOVersion: 0x00000011,

		ImageType:  0x01,
		MiscBitmap: 0x08f7ff7c,

		CanCaptcha: true,
	}
}

func NewClientConfigForAndroidTIM() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForAndroidQiDian() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForAndroidLite() *ClientConfig {
	return &ClientConfig{
		AppID:             0x200300f0,
		CodecAppIDDebug:   "736360370",
		CodecAppIDRelease: "736347652",

		PackageName: "com.tencent.qqlite",
		VersionName: "4.0.2",
		Revision:    "4.0.2.9b6340cd",
		SignatureMD5: []byte{
			0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77,
			0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d,
		},

		BuildTime:  0x0000000060409d2d,
		SDKVersion: "6.0.0.2356",
		SSOVersion: 0x00000005,

		ImageType:  0x01,
		MiscBitmap: 0x00f7ff7c,

		CanCaptcha: true,
	}
}

func NewClientConfigForAndroidTablet() *ClientConfig {
	return &ClientConfig{
		AppID:             0x2002fdd5,
		CodecAppIDDebug:   "73636270;",
		CodecAppIDRelease: "736346857",

		PackageName: "com.tencent.minihd.qq",
		VersionName: "5.9.2",
		Revision:    "5.9.2.3baec0",
		SignatureMD5: []byte{
			0xaa, 0x39, 0x78, 0xf4, 0x1f, 0xd9, 0x6f, 0xf9,
			0x91, 0x4a, 0x66, 0x9e, 0x18, 0x64, 0x74, 0xc7,
		},

		BuildTime:  0x000000005f1e8730,
		SDKVersion: "6.0.0.2433",
		SSOVersion: 0x0000000c,

		ImageType:  0x01,
		MiscBitmap: 0x08f7ff7c,

		CanCaptcha: true,
	}
}

func NewClientConfigForAndroidWatch() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForiOS() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForiOSTIM() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForiOSQiDian() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForiOSLite() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForiPadOS() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForwatchOS() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForWindows() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForWindowsUWP() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForWindowsTIM() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForWindowsTablet() *ClientConfig {
	panic("not implement")
}

func NewClientConfigFormacOS() *ClientConfig {
	panic("not implement")
}

func NewClientConfigForLinux() *ClientConfig {
	panic("not implement")
}
