package config

import (
	"crypto/md5"
	"fmt"
	"log"
	"math/rand"
	"net"

	"github.com/elap5e/go-mobileqq-api/util"
)

type Config struct {
	BaseDir  string
	CacheDir string

	Client *ClientConfig
	Device *DeviceConfig
}

type ClientConfig struct {
	AppID uint32
	FixID uint32

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

	NetworkType  uint16 // 0x0002: Wi-Fi
	NetIPFamily  uint8  // 0x00: Others; 0x01: IPv4; 0x02: IPv6; 0x03: Dual
	IPv4Address  net.IP
	MACAddress   string
	BSSIDAddress string
	SSIDAddress  string

	IMEI string
	IMSI string
	GUID []byte // []byte("%4;7t>;28<fclient.5*6")

	GUIDFlag      uint32
	IsGUIDFileNil bool
	IsGUIDGenSucc bool
	IsGUIDChanged bool
}

func NewClientConfig() *ClientConfig {
	return NewClientConfigForAndroid()
}

func NewClientConfigForAndroid() *ClientConfig {
	return &ClientConfig{
		AppID: 0x20033f65,
		FixID: 0x20033f65, // Release:"7363:367;" Debug:"7363:367;"

		PackageName: "com.tencent.mobileqq",
		VersionName: "8.8.12", // "8.8.12.5675"
		Revision:    "8.8.12.5475818d",
		SignatureMD5: []byte{
			0xa6, 0xb7, 0x45, 0xbf, 0x24, 0xa2, 0xc2, 0x77,
			0x52, 0x77, 0x16, 0xf6, 0xf3, 0x6e, 0xb6, 0x8d,
		},

		BuildTime:  0x0000000060c9bd50,
		SDKVersion: "6.0.0.2477",
		SSOVersion: 0x00000012,

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
		AppID: 0x200300f0,
		FixID: 0x200300f0, // Release:"736360370" Debug:"736360370"

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
		AppID: 0x2002fdd5,
		FixID: 0x2002fdd5, // Release:"736346857" Debug:"73636270;"

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

func NewClientConfigForWindowsQiDian() *ClientConfig {
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

func NewDeviceConfig() *DeviceConfig {
	return NewDeviceConfigForAndroid()
}

func NewDeviceConfigForAndroid() *DeviceConfig {
	mac1 := "00:50:56:C0:00:08"
	osid := "RKQ1.200827.002"
	return &DeviceConfig{
		Bootloader:   "unknown",
		ProcVersion:  "Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)",
		Codename:     "davinci",
		Incremental:  "20.10.20",
		Fingerprint:  "Xiaomi/davinci/davinci:11/" + osid + "/20.10.20:user/release-keys",
		BootID:       "aa6bf49c-a995-4761-874f-8b1a9eee341e",
		OSBuildID:    osid,
		Baseband:     "4.3.c5-00069-SM6150_GEN_PACK-1",
		InnerVersion: "20.10.20",

		NetworkType:  0x0002,
		NetIPFamily:  0x03,
		IPv4Address:  net.IPv4(192, 168, 0, 100),
		MACAddress:   mac1,
		BSSIDAddress: "00:50:56:C0:00:09",
		SSIDAddress:  "unknown",

		IMEI: "860308028836598",
		IMSI: "088906035901507678",
		GUID: util.STBytesTobytes(md5.Sum(append([]byte(osid), mac1...))),

		GUIDFlag:      uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00)),
		IsGUIDFileNil: false,
		IsGUIDGenSucc: true,
		IsGUIDChanged: false,
	}
}

func NewDeviceConfigBySeed(seed int64) *DeviceConfig {
	r := rand.New(rand.NewSource(seed))
	buf := make([]byte, 20)
	_, err := r.Read(buf)
	if err != nil {
		log.Fatalf("failed to generate device config")
	}
	ipv4 := net.IPv4(192, 168, 0, buf[0])
	mac1 := fmt.Sprintf("00:50:%02x:%02x:00:%02x", buf[1], buf[2], buf[0])
	mac2 := fmt.Sprintf("00:50:%02x:%02x:00:%02x", buf[1], buf[2], buf[3])
	uuid := fmt.Sprintf(
		"%02x%02x%02x%02x-%02x%02x-%02x%02x-%02x%02x-%02x%02x%02x%02x%02x%02x",
		buf[4], buf[5], buf[6], buf[7], buf[8], buf[9], buf[10], buf[11],
		buf[12], buf[13], buf[14], buf[15], buf[16], buf[17], buf[18], buf[19],
	)
	imei := fmt.Sprintf("86030802%07d", r.Int31n(10000000))
	imsi := fmt.Sprintf("460001%09d", r.Int31n(1000000000))
	osid := fmt.Sprintf("RKQ1.%07d.002", r.Int31n(10000000))
	return &DeviceConfig{
		Bootloader:   "unknown",
		ProcVersion:  "Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)",
		Codename:     "davinci",
		Incremental:  "20.10.20",
		Fingerprint:  "Xiaomi/davinci/davinci:11/" + osid + "/20.10.20:user/release-keys",
		BootID:       uuid,
		OSBuildID:    osid,
		Baseband:     "4.3.c5-00069-SM6150_GEN_PACK-1",
		InnerVersion: "20.10.20",

		NetworkType:  0x0002,
		NetIPFamily:  0x03,
		IPv4Address:  ipv4,
		MACAddress:   mac1,
		BSSIDAddress: mac2,
		SSIDAddress:  "unknown",

		IMEI: imei,
		IMSI: imsi,
		GUID: util.STBytesTobytes(md5.Sum(append([]byte(osid), mac1...))),

		GUIDFlag:      uint32((1 << 24 & 0xFF000000) | (0 << 8 & 0xFF00)),
		IsGUIDFileNil: false,
		IsGUIDGenSucc: true,
		IsGUIDChanged: false,
	}
}
