package model

type Entry struct {
	Task string
	Start int64
	Duration int64
}

type CurrentEntry struct {
	Task string
	Start int64
}
