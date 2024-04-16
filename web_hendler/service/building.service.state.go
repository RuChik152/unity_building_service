package service

var STATUS_BUILDING bool = false
var STATUS_RESET bool = false
var PID_PROCCES_BUILDING int

type CheckListBuilding struct {
	git       int
	pre_build int
	building  int
	upload    int
}

var CHECK_LIST CheckListBuilding
