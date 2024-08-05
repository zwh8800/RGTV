package embeddata

import (
	_ "embed"
)

const FontName = "fusion-pixel-12px-monospaced.ttc"

//go:embed fusion-pixel-12px-monospaced.ttc
var FontData []byte
