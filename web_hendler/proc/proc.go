package proc

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func getOS() string {
	return runtime.GOOS
}

func GetListChildProcces(pid int) []int {
	procPID := fmt.Sprintf("(ParentProcessId=%d)", pid)
	cmdArgs := []string{
		"process",
		"where",
		procPID,
		"get",
		"ProcessId",
	}
	cmd := exec.Command("wmic", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Ошибка. ", string(output), "err: ", err)
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\r\r\n")
	fmt.Println(lines)
	var pidArray []int

	for index, value := range lines {
		if index > 0 && value != "" {
			fields := strings.Fields(value)
			if len(fields) > 0 {
				number, err := strconv.Atoi(fields[0])
				if err != nil {
					log.Println("ERROR CONVERTING: ", err)
				}

				pidArray = append(pidArray, number)
			}
		}
	}

	return pidArray
}
