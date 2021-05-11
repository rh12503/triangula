package fitness

// fastRound is an optimized version of math.Round.
func fastRound(n float64) int {
	return int(n+0.5) << 0
}

