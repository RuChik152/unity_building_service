package service

import (
	"os/exec"
)

var NAME_APK string
var PROJECT_FOLDER string
var DEST_ANDROID_BUILD_FOLDER = "C:\\Unity\\build\\"
var DEST_WIN_BUILD_FOLDER string
var PATH_TO_LOGS = "C:\\Unity\\logs\\test_runner_logs.log"
var PATH_CURRENT_APK string
var PATH_CURRENT_OBB string
var PATH_TO_EDITOR = "C:\\Program Files\\Unity\\Hub\\Editor\\2022.1.23f1\\Editor\\Unity.exe"
var PROCCES_BUILDING *exec.Cmd

var LIST_PLATFORM map[string][]string = map[string][]string{
	"Android": {"PICO", "OCULUS"},
	"Win64":   {"PC"},
}
