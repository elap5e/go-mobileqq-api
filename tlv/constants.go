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
	defaultDeviceCodename     = []byte("davinci")
	defaultDeviceIncremental  = []byte("20.10.20")
	defaultDeviceFingerprint  = []byte("Xiaomi/davinci/davinci:11/RKQ1.200826.002/20.10.20:user/release-keys")
	defaultDeviceBootID       = []byte("aa6bf49c-a995-4761-874f-8b1a9eee341d")
	defaultDeviceOSBuildID    = []byte("RKQ1.200826.002")
	defaultDeviceBaseband     = []byte("4.3.c5-00069-SM6150_GEN_PACK-1")
	defaultDeviceInnerVersion = []byte("20.10.20")
)
