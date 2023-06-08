package nwc24

var table2 = [8]uint8{0x1, 0x5, 0x0, 0x4, 0x2, 0x3, 0x6, 0x7}
var invertedTable = [16]uint8{0xD, 0x5, 0x9, 0x7, 0x0, 0xF, 0xA, 0x2, 0xC, 0x3, 0xE, 0x1, 0x8, 0x6, 0xB, 0x4}

// WiiNumber includes many of the fields that make up a Wii Number.
type WiiNumber struct {
	// hollywoodId is the ID of the console the Wii number originated from.
	hollywoodId uint32
	// generationCount is equal to the amount of times the Wii console has been restored.
	generationCount uint16
	// hardwareModel is the type of console this number was generated on
	hardwareModel uint8
	// areaCode is the area of the console this number was generated on
	areaCode uint8
	// unscrambled is the unscrambled Wii Number.
	unscrambled uint64
}

func (w *WiiNumber) CheckWiiNumber() bool {
	temp := w.unscrambled
	for i := 0; i <= 42; i++ {
		val := temp >> uint64(52-i)
		if val&1 != 0 {
			val = 0x0000000000000635 << uint64(42-i)
			temp ^= val
		}
	}

	return uint8(temp) == 0
}

func (w *WiiNumber) GetHollywoodID() uint32 {
	return w.hollywoodId
}

func (w *WiiNumber) GetGenerationCount() uint16 {
	return w.generationCount
}

func (w *WiiNumber) GetHardwareModel() uint8 {
	return w.hardwareModel
}

func (w *WiiNumber) GetAreaCode() uint8 {
	return w.areaCode
}

func getByte(value uint64, shift uint8) uint8 {
	return uint8(value >> (shift * 8))
}

func insertByte(value uint64, shift uint8, byte byte) uint64 {
	mask := 0x00000000000000FF << (shift * 8)
	inst := uint64(byte) << (shift * 8)
	return (value & ^uint64(mask)) | inst
}

func unscrambleId(wiiNumber uint64) uint64 {
	wiiNumber &= 0x001FFFFFFFFFFFFF
	wiiNumber ^= 0x00005E5E5E5E5E5E
	wiiNumber &= 0x001FFFFFFFFFFFFF

	mixId := wiiNumber
	mixId ^= 0xFF
	mixId = (wiiNumber << 5) & 0x20

	wiiNumber |= mixId << 48
	wiiNumber >>= 1

	mixId = wiiNumber
	for i := 0; i <= 5; i++ {
		val := getByte(mixId, table2[i])
		wiiNumber = insertByte(wiiNumber, uint8(i), val)
	}

	for i := 0; i <= 5; i++ {
		val := getByte(wiiNumber, uint8(i))
		newByte := ((invertedTable[(val>>4)&0xF]) << 4) | (invertedTable[val&0xF])
		wiiNumber = insertByte(wiiNumber, uint8(i), newByte)
	}

	mixIdCopy := wiiNumber >> 0x20
	anotherMixIdCopy := wiiNumber>>0x16 | (mixIdCopy&0x7FF)<<10
	mixIdCopy = wiiNumber*0x400 | (mixIdCopy >> 0xB & 0x3FF)
	mixIdCopy = (anotherMixIdCopy << 0x20) | mixIdCopy
	mixIdCopy ^= 0x0000B3B3B3B3B3B3

	return mixIdCopy
}

func LoadWiiNumber(wiiNumber uint64) WiiNumber {
	unscrambled := unscrambleId(wiiNumber)
	return WiiNumber{
		hollywoodId:     uint32(unscrambled >> 15),
		generationCount: uint16((unscrambled >> 10) & 0x1F),
		hardwareModel:   uint8((unscrambled >> 47) & 7),
		areaCode:        uint8((unscrambled >> 50) & 7),
		unscrambled:     unscrambled,
	}
}
