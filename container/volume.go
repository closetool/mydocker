package container

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)



func NewWorkSpace(volume, imageName, containerName string) {
	CreateReadOnlyLayer(imageName)
	CreateWriteLayer(containerName)
	CreateMountPoint(containerName,imageName)

	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			MountVolume(volumeURLs,containerName)
			log.Infof("%q",volumeURLs)
		}else{
			log.Errorf("Volume parameter input is not correct.")
		}
	}
}

func CreateReadOnlyLayer(imageName string) error {
	unTarFolderUrl := fmt.Sprintf("%s/%s/",RootUrl,imageName)
	imageUrl := fmt.Sprintf("%s/%s.tar",RootUrl,imageName)
	exist, err := PathExists(unTarFolderUrl)
	if err != nil {
		log.Infof("Fail to judge whether dir %s exists.")
		return err
	}
	if !exist {
		if err := os.MkdirAll(unTarFolderUrl,0777); err != nil {
			log.Errorf("Mkdir dir %s error. %v", unTarFolderUrl,err)
			return err
		}
		if _, err := exec.Command("tar","-xvf",imageUrl,"-C",unTarFolderUrl).CombinedOutput(); err != nil {
			log.Errorf("Untar dir %s error %v", unTarFolderUrl,err)
			return err
		}
	}
	return nil
}

func CreateWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl,containerName)
	if err  := os.MkdirAll(writeURL,0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", writeURL,err)
	}
}

func CreateMountPoint(containerName string, imageName string) error {
	mntUrl := fmt.Sprintf(MntUrl,containerName)
	if err := os.MkdirAll(mntUrl,0777); err != nil {
		log.Errorf("Mkdir dir %s error. %v", mntUrl,err)
		return err
	}
	tmpWriteLayer := fmt.Sprintf(WriteLayerUrl,containerName)
	tmpImageLocation := RootUrl + "/" + imageName
	mntURL := fmt.Sprintf(MntUrl, containerName)
	dirs := "dirs="+tmpWriteLayer+":" + tmpImageLocation
	_, err := exec.Command("mount", "-t","aufs","-o",dirs,"none",mntURL).CombinedOutput()
	if err != nil {
		log.Errorf("Run command for creating mount point failed %v",err)
		return err
	}
	return nil
}

func DeleteWorkSpace(volume,containerName string) {
	if volume != "" {
		volumeURLs := strings.Split(volume,":")
		length:= len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {
			DeleteMountPointWithVolume(volumeURLs,containerName)
		}else{
			DeleteMountPoint(containerName)
		}
	}else{
		DeleteMountPoint(containerName)
	}
	DeleteWriteLayer(containerName)
}

func DeleteMountPoint(containerName string) error {
	mntURL := fmt.Sprintf(MntUrl,containerName)
	_, err := exec.Command("umount",mntURL).CombinedOutput()
	if err != nil {
		log.Errorf("Umount %s error %v",mntURL,err)
		return err
	}
	if err := os.RemoveAll(mntURL); err != nil {
		log.Errorf("Remove mountpoint dir %s error %v", mntURL,err)
		return err
	}
	return nil
}

func DeleteMountPointWithVolume(volumeURLs []string,containerName string) error {
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerUrl := mntURL + "/" + volumeURLs[1]
	if _, err := exec.Command("umount", containerUrl).CombinedOutput(); err != nil {
		log.Errorf("Umount volume %s failed. %v", containerUrl,err)
		return err
	}
	if _, err := exec.Command("umount", mntURL).CombinedOutput(); err != nil {
		log.Errorf("Umount mountpoint %s failed. %v",mntURL,err)
		return err
	}
	
	if err := os.RemoveAll(mntURL); err != nil {
		log.Error("Remove mountpoint dir %s error %v", mntURL, err)
		return err
	}
	return nil
}

func DeleteWriteLayer(containerName string) {
	writeURL := fmt.Sprintf(WriteLayerUrl,containerName)
	if err := os.RemoveAll(writeURL); err != nil {
		log.Errorf("Remove dir %s error %v", writeURL, err)
	}
}

func MountVolume(volumeURLs []string, containerName string) error {
	parentUrl := volumeURLs[0]
	if err := os.MkdirAll(parentUrl,0777); err != nil {
		log.Infof("Mkdir parent dir %s error. %v",parentUrl,err)
	}
	containerUrl := volumeURLs[1]
	mntURL := fmt.Sprintf(MntUrl, containerName)
	containerVolumeURL := mntURL +"/"+ containerUrl
	if err := os.MkdirAll(containerVolumeURL,0777); err != nil {
		log.Infof("Mkdir container dir %s error. %v",containerVolumeURL,err)
		return err
	}
	dirs := "dirs=" + parentUrl
	_, err := exec.Command("mount","-t","aufs","-o",dirs,"none",containerVolumeURL).CombinedOutput()
	if err != nil {
		log.Errorf("Mount volume failed. %v",err)
		return err
	}
	return nil
}