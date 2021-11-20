package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	log "github.com/sirupsen/logrus"
)

var (
	RootUrl = "/root"
	MntUrl = "/root/mnt/%s"
	WriteLayerUrl = "/root/writeLayer/%s"
)

var (
	RUNNING = "running"
	STOP = "stopped"
	Exit = "exited"
	DefaultInfoLocation = "/var/run/mydocker/%s/"
	ConfigName = "config.json"
	ContainerLogFile = "container.log"
)

type ContainerInfo struct {
	Pid string `json:"pid"`
	Id string `json:"id"`
	Name string `json:"name"`
	Command string `json:"command"`
	CreateTime string `json:"createTime"`
	Status string `json:"status"`
	Volume string `json:"volume"`
}

func NewParentProcess(tty bool,containerName,volume,imageName string,envSlice []string) (*exec.Cmd, *os.File) {
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		log.Errorf("New pipe error %v", err)
		return nil, nil
	}
	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}else{
		dirURL := fmt.Sprintf(DefaultInfoLocation,containerName)
		if err := os.MkdirAll(dirURL,0622); err != nil {
			log.Errorf("NewParentProcess mkdir %s error %v",dirURL, err)
			return nil,nil
		}
		stdLogFilePath := dirURL + ContainerLogFile
		stdLogFile, err := os.Create(stdLogFilePath)
		if err != nil {
			log.Errorf("NewParentProcess create file %s error %v",stdLogFilePath,err)
			return nil, nil
		}
		cmd.Stdout = stdLogFile
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	cmd.Env = append(os.Environ(),envSlice...)
	// 创建联合文件系统mnt目录
	// 并将程序的工作路径设置为mnt
	// 实现了工作路径上的虚拟
	NewWorkSpace(volume, imageName, containerName)
	cmd.Dir = fmt.Sprintf(MntUrl,containerName)
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}



func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func volumeUrlExtract(volume string) []string {
	return strings.Split(volume,":")
}

