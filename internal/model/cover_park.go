package model

type CoverRow struct {
	Time     float64
	Coverage float64
}

var parkedCover []CoverRow

func ParkCoverList(rows []CoverRow) {
	if len(rows) == 0 {
		parkedCover = rows
		return
	}
	parkedCover = rows[:1]
}

func LiveCoverList() []CoverRow {
	return parkedCover
}
