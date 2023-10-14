package server

import (
	"math"
	"strconv"
)

type PageInfo struct {
	Number int
}

func genPageInfo(total, perPage int64) []PageInfo {
	var ps []PageInfo

	pagesTotal := math.Ceil(float64(total) / float64(perPage))

	for i := 0; i < int(pagesTotal); i++ {
		ps = append(ps, PageInfo{
			Number: i + 1,
		})
	}

	return ps
}

func parseOptionalIntQueryParam(p string, d int) (int, error) {
	if p == "" {
		return d, nil
	}

	result, err := strconv.Atoi(p)
	if err != nil {
		return 0, err
	}

	return result, nil
}
