package server

import (
	"math"

	"alinea.com/internal/service"
)

type PageInfo struct {
	Number int
}

func genPageInfo(p service.PageJson) []PageInfo {
	var ps []PageInfo

	pagesTotal := math.Ceil(float64(p.Total) / float64(p.PerPage))

	for i := 0; i < int(pagesTotal); i++ {
		ps = append(ps, PageInfo{
			Number: i + 1,
		})
	}

	return ps
}
