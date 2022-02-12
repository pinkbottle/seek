package seekcli

import "regexp"

var exp = regexp.MustCompile(`\[([0-9]*):([0-9]*)\]`)

func parsePagination(query string) (int, int) {
	p := exp.FindStringSubmatch(query)
	return int(p[1][0] - '0'), int(p[2][0] - '0')
}
