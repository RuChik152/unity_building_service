package service

import (
	"os/exec"
)

var NAME_KEYSTORE string
var NAME_APK string
var PROJECT_FOLDER string
var DEST_ANDROID_BUILD_FOLDER string
var DEST_WIN_BUILD_FOLDER string
var PATH_TO_LOGS string
var PATH_CURRENT_APK string
var PATH_CURRENT_OBB string
var PATH_TO_EDITOR string
var PROCCES_BUILDING *exec.Cmd
var PATH_BUILDER_MOD string
var PATH_CLOUD_MOD string
var PATH_TO_CONFIG_JSON string

var MAP_CONFIG_DATA map[string]string

var LIST_PLATFORM map[string][]string = map[string][]string{
	//"Android": {"OCULUS", "PICO"},
	"Win64": {"PC"},
}

var PICO_APP_ID string
var PICO_APP_SECRET string
var OCULUS_APP_ID string
var OCULUS_APP_SECRET string
var STEAM_LOGIN string
var STEAM_PASS string
