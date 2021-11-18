package stdgomods

func IntPct(progress, total int64) int64 {
	if total == 0 {
		return 100
	}
	pct := progress * 100 / total
	return pct
}
