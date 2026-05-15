package conv

func ptrOrZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
