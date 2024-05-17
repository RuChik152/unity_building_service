package proc

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"web_hendler/loger"

	"github.com/shirou/gopsutil/process"
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
		loger.LogPrint.Package("PROC").Log(fmt.Sprint("Ошибка. ", string(output), "err: ", err))
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
					loger.LogPrint.Package("PROC").Log(fmt.Sprint("ERROR CONVERTING: ", err))
				}

				pidArray = append(pidArray, number)
			}
		}
	}

	return pidArray
}

func DestroyedBuilding(pid int) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		loger.LogPrint.Package("PROC").Log(fmt.Sprint("Ошибка при получении процесса:", err))
		return
	}

	children, err := p.Children()
	if err != nil {
		loger.LogPrint.Package("PROC").Log(fmt.Sprint("Ошибка при получении дочерних процессов:", err))
		return
	}

	for _, child := range children {
		if err := child.Terminate(); err != nil {
			loger.LogPrint.Package("PROC").Log(fmt.Sprint("Ошибка при завершении дочернего процесса:", err))
		}
	}

	if err := p.Terminate(); err != nil {
		loger.LogPrint.Package("PROC").Log(fmt.Sprint("Ошибка при завершении родительского процесса:", err))
	}
}
