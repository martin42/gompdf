package style

type HAlign string

const (
	HAlignLeft   HAlign = "left"
	HAlignRight  HAlign = "right"
	HAlignCenter HAlign = "center"
)

type VAlign string

const (
	VAlignTop    VAlign = "top"
	VAlignMiddle VAlign = "middle"
	VAlignBottom VAlign = "bottom"
)

type Align struct {
	HAlign `style:"h-align"`
	VAlign `style:"v-align"`
}
