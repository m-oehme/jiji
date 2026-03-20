package common

import "testing"

func TestSplitHorizontal(t *testing.T) {
	tests := []struct {
		name        string
		totalW      int
		leftPercent int
		wantLeft    int
		wantRight   int
	}{
		{"30/70 split on 100", 100, 30, 30, 70},
		{"50/50 split on 100", 100, 50, 50, 50},
		{"30/70 split on 80", 80, 30, 24, 56},
		{"1% left on 100", 100, 1, 1, 99},
		{"99% left on 100", 100, 99, 99, 1},
		{"very small terminal", 3, 30, 1, 2},
		{"width of 1", 1, 50, 0, 1},
		{"width of 2", 2, 50, 1, 1},
		{"zero width", 0, 30, 0, 0},
		{"negative width", -5, 30, 0, 0},
		{"percent below 1 clamped", 100, 0, 1, 99},
		{"percent above 99 clamped", 100, 100, 99, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			left, right := SplitHorizontal(tt.totalW, tt.leftPercent)
			if left != tt.wantLeft || right != tt.wantRight {
				t.Errorf("SplitHorizontal(%d, %d) = (%d, %d), want (%d, %d)",
					tt.totalW, tt.leftPercent, left, right, tt.wantLeft, tt.wantRight)
			}
			// Invariant: left + right == totalW (when totalW > 0)
			if tt.totalW > 0 && left+right != tt.totalW {
				t.Errorf("left(%d) + right(%d) = %d, want %d", left, right, left+right, tt.totalW)
			}
		})
	}
}

func TestInnerSize(t *testing.T) {
	tests := []struct {
		name      string
		outerW    int
		outerH    int
		hasBorder bool
		wantW     int
		wantH     int
	}{
		{"with border", 80, 24, true, 78, 22},
		{"without border", 80, 24, false, 80, 24},
		{"tiny with border", 2, 2, true, 0, 0},
		{"tiny without border", 2, 2, false, 2, 2},
		{"1x1 with border", 1, 1, true, 0, 0},
		{"0x0 with border", 0, 0, true, 0, 0},
		{"0x0 without border", 0, 0, false, 0, 0},
		{"negative with border", -5, -3, true, 0, 0},
		{"negative without border", -5, -3, false, 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, h := InnerSize(tt.outerW, tt.outerH, tt.hasBorder)
			if w != tt.wantW || h != tt.wantH {
				t.Errorf("InnerSize(%d, %d, %v) = (%d, %d), want (%d, %d)",
					tt.outerW, tt.outerH, tt.hasBorder, w, h, tt.wantW, tt.wantH)
			}
		})
	}
}
