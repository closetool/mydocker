package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"

	"github.com/closetool/mydocker/container"
	log "github.com/sirupsen/logrus"
)


func stopContainer(containerName string) {
	pidStr, err := getContainerPidByName(containerName) 
	if err != nil {
		log.Fatalf("Get container pid by name %s error %v",containerName,err)
	}
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		log.Fatalf("Convert pid from string to int error %v",err)
	}
	if err := syscall.Kill(pid,syscall.SIGTERM); err != nil {
		log.Fatalf("Stop container %s error %v",containerName, err)
	}
	containerInfo, err := getContainerInfoByName(containerName)
	if err != nil {
		log.Fatalf("Get container %s info error %v",containerName, err)
	}
	containerInfo.Status = container.STOP
	containerInfo.Pid = " "
	newContentBytes, err := json.Marshal(containerInfo)
	if err != nil {
		log.Fatalf("Json marshal %s error %v", containerName, err)
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation,containerName)
	configFilePath := dirURL + container.ConfigName
	if err := ioutil.WriteFile(configFilePath,newContentBytes,0622); err != nil {
		log.Errorf("Write file %s error %v", configFilePath ,err)
	}
}

func getContainerInfoByName(containerName string)(*container.ContainerInfo,error) {
	dirURL := fmt.Sprintf(container.DefaultInfoLocation,containerName)
	configFilePath := dirURL + container.ConfigName
	contentBytes, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		log.Errorf("Read file %s error %v",configFilePath, err)
		return nil ,err
	}
	var containerInfo container.ContainerInfo
	if err := json.Unmarshal(contentBytes,&containerInfo); err != nil {
		log.Errorf("GetContainerInfoByName unmarshal error %v",err)
		return nil, err
	}
	return &containerInfo, nil
}

func removeContainer(containerName string) {
	containerInfo,err := getContainerInfoByName(containerName)
	if err != nil {
		log.Fatal("Get container %s info error %v",containerName,err)
	}
	if containerInfo.Status != container.STOP {
		log.Fatal("Couldn't remove running container")
	}
	dirURL := fmt.Sprintf(container.DefaultInfoLocation,containerName)
	if err := os.RemoveAll(dirURL); err != nil {
		log.Fatal("Remove file %s error %v",dirURL, err)
	}
}