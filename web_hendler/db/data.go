package db

type CommitData struct {
	ID      int
	AUTHOR  string
	COMMENT string
	SHA     string
}

type GetCommit struct {
	ID      int
	AUTHOR  string
	MESSAGE string
	SHA     string
}

type CheckData struct {
	ID int
}

var Commit CommitData
