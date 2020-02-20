package coalbox

import "image"

type BoundingBox struct {
	Text                                                 string
	TopLeftX, TopLeftY, TopRightX, TopRightY             int
	BottomLeftX, BottomLeftY, BottomRightX, BottomRightY int
	Confidence                                           float64
}

func BoxFromRectangle(r image.Rectangle) BoundingBox {
	return BoundingBox{
		TopLeftX:     r.Min.X,
		TopLeftY:     r.Min.Y,
		TopRightX:    r.Max.X,
		TopRightY:    r.Min.Y,
		BottomLeftX:  r.Min.X,
		BottomLeftY:  r.Max.Y,
		BottomRightX: r.Max.X,
		BottomRightY: r.Max.Y,
	}
}

func BoxesFromRectangles(rs []image.Rectangle) []BoundingBox {
	bbs := make([]BoundingBox, len(rs), len(rs))
	for i, _ := range rs {
		bbs[i] = BoxFromRectangle(rs[i])
	}
	return bbs
}
