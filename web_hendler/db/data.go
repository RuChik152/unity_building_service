package db

type CommitData struct {
	ID      int
	AUTHOR  string
	COMMENT string
	SHA     string
}

type CheckData struct {
	ID int
}

var Commit CommitData
