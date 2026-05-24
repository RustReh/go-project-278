package handler

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/RustReh/go-project-278/internal/apperr"
)

var rangeQueryRE = regexp.MustCompile(`^\[\s*(\d+)\s*,\s*(\d+)\s*\]$`)

// parseRangeQuery разбирает range=[start,end] (полуинтервал [start, end)).
func parseRangeQuery(raw string) (start, end int, err error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0, 0, apperr.Validation("range query parameter is required", map[string]any{"range": raw})
	}

	m := rangeQueryRE.FindStringSubmatch(raw)
	if m == nil {
		return 0, 0, apperr.Validation("invalid range format, expected [start,end]", map[string]any{"range": raw})
	}

	start, err = strconv.Atoi(m[1])
	if err != nil || start < 0 {
		return 0, 0, apperr.Validation("invalid range start", map[string]any{"range": raw})
	}
	end, err = strconv.Atoi(m[2])
	if err != nil || end < start {
		return 0, 0, apperr.Validation("invalid range end", map[string]any{"range": raw})
	}

	return start, end, nil
}

func contentRangeHeader(start, end int, total int64) string {
	return "links " + strconv.Itoa(start) + "-" + strconv.Itoa(end) + "/" + strconv.FormatInt(total, 10)
}
