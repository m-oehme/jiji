package common

// SplitHorizontal divides totalWidth into left and right widths
// based on leftPercent (1-99). Remainder goes to the right pane.
func SplitHorizontal(totalWidth, leftPercent int) (leftW, rightW int) {
	if totalWidth <= 0 {
		return 0, 0
	}
	if leftPercent < 1 {
		leftPercent = 1
	}
	if leftPercent > 99 {
		leftPercent = 99
	}
	leftW = totalWidth * leftPercent / 100
	if leftW < 1 {
		leftW = 1
	}
	rightW = totalWidth - leftW
	if rightW < 1 {
		rightW = 1
		leftW = totalWidth - 1
	}
	if leftW < 1 {
		leftW = 0
		rightW = totalWidth
	}
	return leftW, rightW
}

// InnerSize computes the usable interior dimensions after subtracting
// border chrome (1 cell per side for border).
func InnerSize(outerW, outerH int, hasBorder bool) (innerW, innerH int) {
	if !hasBorder {
		return max(outerW, 0), max(outerH, 0)
	}
	innerW = outerW - 2 // left + right border
	innerH = outerH - 2 // top + bottom border
	if innerW < 0 {
		innerW = 0
	}
	if innerH < 0 {
		innerH = 0
	}
	return innerW, innerH
}
