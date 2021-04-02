package draw

import (
	"fmt"
	"testing"
)

func TestQuarterEllipseOutlines(t *testing.T) {
	// empty ellipses
	quarter(t, 0, 0)
	quarter(t, 0, 1)
	quarter(t, 1, 0)
	// circles
	quarter(t,
		1, 1,
		0, 0,
	)
	quarter(t,
		2, 2,
		0, 0,
	)
	quarter(t,
		3, 3,
		0, 1, 1, 0,
	)
	quarter(t,
		4, 4,
		0, 1, 1, 0,
	)
	quarter(t,
		5, 5,
		0, 2, 1, 2, 2, 1, 2, 0,
	)
	quarter(t,
		6, 6,
		0, 2, 1, 2, 2, 1, 2, 0,
	)
	quarter(t,
		7, 7,
		0, 3, 1, 3, 2, 2, 3, 1, 3, 0,
	)
	quarter(t,
		8, 8,
		0, 3, 1, 3, 2, 2, 3, 1, 3, 0,
	)
	quarter(t,
		9, 9,
		0, 4, 1, 4, 2, 3, 3, 3, 3, 2, 4, 1, 4, 0,
	)
	quarter(t,
		10, 10,
		0, 4, 1, 4, 2, 3, 3, 3, 3, 2, 4, 1, 4, 0,
	)
	// asymmetric ellipses
	quarter(t,
		2, 1,
		0, 0,
	)
	quarter(t,
		1, 2,
		0, 0,
	)
	quarter(t,
		3, 1,
		0, 0, 1, 0,
	)
	quarter(t,
		1, 3,
		0, 1, 0, 0,
	)
}

func TestEllipseOutlines(t *testing.T) {
	outline(t, 2, 3, 0, 0)
	outline(t,
		2, 3, 3, 3,
		3, 5, 4, 4, 3, 3, 2, 4,
	)
	outline(t,
		2, 3, 3, 4,
		3, 6, 4, 5, 4, 4, 3, 3, 2, 4, 2, 5,
	)
	outline(t,
		2, 3, 4, 3,
		4, 5, 5, 4, 4, 3, 3, 3, 2, 4, 3, 5,
	)
	outline(t,
		4, 3, 4, 4,
		6, 6, 7, 5, 7, 4, 6, 3, 5, 3, 4, 4, 4, 5, 5, 6,
	)
	outline(t,
		3, 2, 5, 5,
		5, 6, 6, 6, 7, 5, 7, 4, 7, 3, 6, 2, 5, 2, 4, 2, 3, 3, 3, 4, 3, 5, 4, 6,
	)
}

func TestEllipseArea(t *testing.T) {
	area(t, 2, 3, 0, 0)
	area(t,
		2, 3, 3, 3,
		3, 5, 3, 5, 3, 3, 3, 3, 2, 4, 4, 4,
	)
	area(t,
		2, 3, 3, 4,
		3, 6, 3, 6, 3, 3, 3, 3, 2, 5, 4, 5, 2, 4, 4, 4,
	)
	area(t,
		2, 3, 4, 3,
		3, 5, 4, 5, 3, 3, 4, 3, 2, 4, 5, 4,
	)
	area(t,
		4, 3, 4, 4,
		5, 6, 6, 6, 5, 3, 6, 3, 4, 5, 7, 5, 4, 4, 7, 4,
	)
	area(t,
		3, 2, 5, 5,
		4, 6, 6, 6, 4, 2, 6, 2, 3, 5, 7, 5, 3, 3, 7, 3, 3, 4, 7, 4,
	)
}

func quarter(t *testing.T, w, h int, wantXYs ...int) {
	check(t,
		fmt.Sprintf("quarter ellipse(%v,%v)", w, h),
		quaterEllipsePoints(w, h),
		wantXYs...,
	)
}

func outline(t *testing.T, x, y, w, h int, wantXYs ...int) {
	check(t,
		fmt.Sprintf("ellipse outline(%v,%v,%v,%v)", x, y, w, h),
		ellipseOutline(x, y, w, h),
		wantXYs...,
	)
}

func area(t *testing.T, x, y, w, h int, wantXYs ...int) {
	check(t,
		fmt.Sprintf("ellipse area(%v,%v,%v,%v)", x, y, w, h),
		ellipseArea(x, y, w, h),
		wantXYs...,
	)
}

func check(t *testing.T, header string, got []point, wantXYs ...int) {
	want := make([]point, len(wantXYs)/2)
	for i := range want {
		want[i].x = wantXYs[i*2]
		want[i].y = wantXYs[i*2+1]
	}
	fail := func() {
		t.Errorf(
			"%v\nwant\n%v\ngot\n%v",
			header, want, got,
		)
	}
	if len(got) != len(want) {
		fail()
		return
	}
	for i := range got {
		if got[i] != want[i] {
			fail()
			return
		}
	}
}
