package draw

// runeToFont maps a unicode rune to the index of the respective glyph in the
// font bitmap. The bitmap contains only a subset of all existing runes, if r is
// not present in the bitmap, a replacement character is returned.
func runeToFont(r rune) rune {
	if 0 <= r && r <= 127 {
		return r
	}
	from128 := []rune("ÇüéâäàåçêëèÏÎÌÄÅÈæÖöÜß§²³")
	for i, s := range from128 {
		if r == s {
			return rune(128 + i)
		}
	}
	return 0 // character not found
}
