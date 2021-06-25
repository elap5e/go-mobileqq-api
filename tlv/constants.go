package tlv

func SetSSOVersion(ver uint32) {
	defaultSSOVersion = ver
}

func SetDeviceOSBuildID(id []byte) {
	defaultDeviceOSBuildID = id
}

var (
	defaultSSOVersion = uint32(0x00000011)

	defaultDeviceBootloader   = []byte("unknown")
	defaultDeviceProcVersion  = []byte("Linux version 2.6.18-92.el5 (brewbuilder@ls20-bc2-13.build.redhat.com)")
	defaultDeviceCodename     = []byte("davinci")
	defaultDeviceIncremental  = []byte("20.10.20")
	defaultDeviceFingerprint  = []byte("Xiaomi/davinci/davinci:11/RKQ1.200827.002/20.10.20:user/release-keys")
	defaultDeviceBootID       = []byte("aa6bf49c-a995-4761-874f-8b1a9eee341e")
	defaultDeviceOSBuildID    = []byte("RKQ1.200827.002")
	defaultDeviceBaseband     = []byte("4.3.c5-00069-SM6150_GEN_PACK-1")
	defaultDeviceInnerVersion = []byte("20.10.20")
)
