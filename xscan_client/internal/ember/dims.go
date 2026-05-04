package ember

// Block dimensions aligned with go_server/ember_cert/features.py PEFeatureExtractor order.
const (
	DimGeneral      = GeneralFileInfoDim // 7
	DimHistogram    = ByteHistogramDim   // 256
	DimByteEntropy  = ByteEntropyHistogramDim
	DimStrings      = 177
	DimHeader       = 74
	DimSection      = 224
	DimImports      = 1282 // 2 + 256 + 1024
	DimExports      = 129  // 1 + 128
	DimDataDir      = 34   // 16*2 + 2
	DimRich         = 33   // 1 + 32
	DimAuthenticode = 8
	DimWarnings     = 88 // 87 + 1
)

// FeatureDimFull is the concatenated Ember PE vector length (~2568).
const FeatureDimFull = DimGeneral + DimHistogram + DimByteEntropy + DimStrings + DimHeader +
	DimSection + DimImports + DimExports + DimDataDir + DimRich + DimAuthenticode + DimWarnings

const (
	offGeneral = 0
	offHist    = offGeneral + DimGeneral
	offByteEnt = offHist + DimHistogram
	offStrings = offByteEnt + DimByteEntropy
	offHeader  = offStrings + DimStrings
	offSection = offHeader + DimHeader
	offImports = offSection + DimSection
	offExports = offImports + DimImports
	offDataDir = offExports + DimExports
	offRich    = offDataDir + DimDataDir
	offAuth    = offRich + DimRich
	offWarn    = offAuth + DimAuthenticode
)
