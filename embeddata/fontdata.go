package embeddata

import (
	_ "embed"
)

const FontName = "ark-pixel-12px-monospaced.ttc"

//go:embed ark-pixel-12px-monospaced.ttc
var FontData []byte
