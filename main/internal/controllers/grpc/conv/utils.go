package conv

import "time"

func parseProtoTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)
	return t
}

func ptrOrZero(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

func intToPtr(i int) *int {
	if i == 0 {
		return nil
	}
	res := i
	return &res
}
