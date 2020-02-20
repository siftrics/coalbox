package coalbox

import (
	"math"
	"sort"
)

func ToSentences(bbs []BoundingBox) []BoundingBox {
	return ToSentencesUsingRatios(bbs, 1.0, 0.5)
}

type taggedBoundingBox struct {
	BoundingBox
	tag int
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func ToSentencesUsingRatios(bbs []BoundingBox, spaceRatio, overlapRatio float64) []BoundingBox {
	taggedBbs := make([]taggedBoundingBox, len(bbs), len(bbs))
	for i, _ := range bbs {
		taggedBbs[i].BoundingBox = bbs[i]
		taggedBbs[i].tag = -1
	}
	sort.Slice(
		taggedBbs,
		func(i, j int) bool {
			return min(taggedBbs[i].TopLeftX, taggedBbs[i].BottomLeftX) < min(taggedBbs[j].TopLeftX, taggedBbs[j].BottomLeftX)
		},
	)
	tagFountain := 0
	numTagged := 0
	for numTagged < len(taggedBbs) {
		i := 0
		for ; i < len(taggedBbs); i++ {
			if taggedBbs[i].tag == -1 {
				break
			}
		}
		taggedBbs[i].tag = tagFountain
		numTagged++
		maxTopY := max(taggedBbs[i].TopLeftY, taggedBbs[i].TopRightY)
		minTopY := min(taggedBbs[i].TopLeftY, taggedBbs[i].TopRightY)
		maxBottomY := max(taggedBbs[i].BottomLeftY, taggedBbs[i].BottomRightY)
		minBottomY := min(taggedBbs[i].BottomLeftY, taggedBbs[i].BottomRightY)
		height := float64(maxBottomY - minTopY)
		avgWidth := float64(max(taggedBbs[i].TopRightX, taggedBbs[i].BottomRightX) - min(taggedBbs[i].TopLeftX, taggedBbs[i].BottomLeftX))
		n := 1.0
		for j, _ := range taggedBbs {
			if i == j || taggedBbs[j].tag != -1 {
				continue
			}

			if (maxTopY >= min(taggedBbs[j].TopLeftY, taggedBbs[j].TopRightY) && // overlap
				minTopY <= max(taggedBbs[j].BottomLeftY, taggedBbs[j].BottomRightY) &&
				float64(max(taggedBbs[j].BottomLeftY, taggedBbs[j].BottomRightY)-minTopY)/height > overlapRatio) ||
				(maxBottomY >= min(taggedBbs[j].TopLeftY, taggedBbs[j].TopRightY) && // overlap
					minBottomY <= max(taggedBbs[j].BottomLeftY, taggedBbs[j].BottomRightY) &&
					float64(maxBottomY-min(taggedBbs[j].TopLeftY, taggedBbs[j].TopRightY))/height > overlapRatio) ||
				(minTopY <= max(taggedBbs[j].TopLeftY, taggedBbs[j].TopRightY) && // fully contained
					maxBottomY >= min(taggedBbs[j].BottomLeftY, taggedBbs[j].BottomRightY)) {

				spaceBetweenThisBbAndPrevBb := min(taggedBbs[j].TopLeftX, taggedBbs[j].BottomLeftX) - max(taggedBbs[i].TopRightX, taggedBbs[i].BottomRightX)
				if float64(spaceBetweenThisBbAndPrevBb)/avgWidth > spaceRatio {
					break
				}
				taggedBbs[j].tag = tagFountain
				numTagged++
				i = j
				width := float64(max(taggedBbs[j].TopRightX, taggedBbs[j].BottomRightX) - min(taggedBbs[j].TopLeftX, taggedBbs[j].BottomLeftX))
				avgWidth = ((avgWidth * n) + width) / (n + 1.0)
				n += 1.0
			}
		}
		tagFountain++
	}
	tag2TopY := make(map[int]int)
	for i, _ := range taggedBbs {
		if _, ok := tag2TopY[taggedBbs[i].tag]; ok {
			continue
		}
		topY := min(taggedBbs[i].TopLeftY, taggedBbs[i].TopRightY)
		for j, _ := range taggedBbs {
			if taggedBbs[j].tag == taggedBbs[i].tag &&
				min(taggedBbs[j].TopLeftY, taggedBbs[j].TopRightY) < topY {
				topY = min(taggedBbs[j].TopLeftY, taggedBbs[j].TopRightY)
			}
		}
		tag2TopY[taggedBbs[i].tag] = topY
	}
	sort.Slice(
		taggedBbs,
		func(i, j int) bool {
			if taggedBbs[i].tag == taggedBbs[j].tag {
				return min(taggedBbs[i].TopLeftX, taggedBbs[i].BottomLeftX) < min(taggedBbs[j].TopLeftX, taggedBbs[j].BottomLeftX)
			}
			return tag2TopY[taggedBbs[i].tag] < tag2TopY[taggedBbs[j].tag]
		},
	)
	sentences := make([]BoundingBox, len(tag2TopY), len(tag2TopY))
	if len(sentences) == 0 {
		return sentences
	}
	tag := taggedBbs[0].tag
	minX := min(taggedBbs[0].TopLeftX, taggedBbs[0].BottomLeftX)
	minY := min(taggedBbs[0].TopLeftY, taggedBbs[0].TopRightY)
	maxX, maxY := math.MinInt32, math.MinInt32
	text := ""
	confidence := 0.0
	n := 0.0
	for _, bb := range taggedBbs {
		if bb.tag != tag {
			sentences[tag].Text = text
			sentences[tag].TopLeftX = minX
			sentences[tag].TopLeftY = minY
			sentences[tag].TopRightX = maxX
			sentences[tag].TopRightY = minY
			sentences[tag].BottomLeftX = minX
			sentences[tag].BottomLeftY = maxY
			sentences[tag].BottomRightX = maxX
			sentences[tag].BottomRightY = maxY
			sentences[tag].Confidence = confidence
			minX = min(bb.TopLeftX, bb.BottomLeftX)
			minY = min(bb.TopLeftY, bb.TopRightY)
			maxX, maxY = math.MinInt32, math.MinInt32
			text = ""
			confidence = 0.0
			n = 0.0
			tag = bb.tag
		}
		maxX = max(maxX, max(bb.TopRightX, bb.BottomRightX))
		maxY = max(maxY, max(bb.BottomLeftY, bb.BottomRightY))
		minY = min(minY, min(bb.TopLeftY, bb.TopRightY))
		if text == "" {
			text += bb.Text
		} else {
			text += " " + bb.Text
		}
		confidence = (confidence*n + bb.Confidence) / (n + 1.0)
		n += 1.0
	}
	sentences[tag].Text = text
	sentences[tag].TopLeftX = minX
	sentences[tag].TopLeftY = minY
	sentences[tag].TopRightX = maxX
	sentences[tag].TopRightY = minY
	sentences[tag].BottomLeftX = minX
	sentences[tag].BottomLeftY = maxY
	sentences[tag].BottomRightX = maxX
	sentences[tag].BottomRightY = maxY
	sentences[tag].Confidence = confidence
	return sentences
}
