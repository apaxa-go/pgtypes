package uuid

const uuidLen = 16                     // UUID length in bytes
const uuidCleanStringLen = uuidLen * 2 // Raw UUID string representation length in chars ("6ba7b8149dad11d180b400c04fd430c8")
const uuidStringLen = uuidLen*2 + 4    // Default UUID string representation length in chars ("6ba7b814-9dad-11d1-80b4-00c04fd430c8")

// UUID parts lengths in bytes
const (
	part0Len = 4
	part1Len = 2
	part2Len = 2
	part3Len = 2
	part4Len = 6
)

const uuidDelim = '-'

// UUID delimiters positions in standard string representation
const (
	sDelim0At = part0Len * 2
	sDelim1At = sDelim0At + 1 + part1Len*2
	sDelim2At = sDelim1At + 1 + part2Len*2
	sDelim3At = sDelim2At + 1 + part3Len*2
)

// UUID parts position in standard string representation
const (
	sPart0From = 0
	sPart0To   = sPart0From + part0Len*2
	sPart1From = sPart0To + 1
	sPart1To   = sPart1From + part1Len*2
	sPart2From = sPart1To + 1
	sPart2To   = sPart2From + part2Len*2
	sPart3From = sPart2To + 1
	sPart3To   = sPart3From + part3Len*2
	sPart4From = sPart3To + 1
	sPart4To   = sPart4From + part4Len*2
)

// UUID parts position in internal representation
// No need to store partXTo as partXTo = part(X+1)From
const (
	part0From = 0
	part1From = part0From + part0Len
	part2From = part1From + part1Len
	part3From = part2From + part2Len
	part4From = part3From + part3Len
)
