package commonTestHelper

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	uiTest "github.com/easysoft/zentaoatf/test/helper/zentao/ui"
)

func Run(version string, codeDir string) (err error) {
	versionNumber := strings.ReplaceAll(version, ".", "")
	// codeDir = "/www/zentaopms" + versionNumber

	_, err = os.Stat(codeDir)
	if os.IsExist(err) {
		os.RemoveAll(codeDir)
	}

	//docker run --name zentao -p 8081:80 --network=zentaonet -v D:\docker\www\zentaopms:/www/zentaopms -v D:\docker\mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d easysoft/zentao:12.3.3
	// cmd := exec.Command("docker", "run", "--name", "zentao"+versionNumber, "-p", "8081:80", "-v", codeDir+":/www/zentaopms", "-e", "-d", "easysoft/zentao:"+version)
	cmd := exec.Command("docker", "run", "--name", "zentao"+versionNumber, "-p", "8081:80", "-v", codeDir+":/www/zentaopms", "-d", "easysoft/zentao:"+version)
	fmt.Println(cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		return
	}
	if strings.Contains(string(output), "exit") {
		return errors.New("Run docker fail")
	}
	return err
}

func IsExistContainer(name string) bool {
	cmd := exec.Command("docker", "ps", "-a", "--format", "'{{.Names}}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), name)
}

func IsRuning(name string) bool {
	cmd := exec.Command("docker", "ps", "--format", "'{{.Names}}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), name)
}

func List() []string {
	cmd := exec.Command("docker", "ps", "--format", "'{{.Names}}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil
	}
	dockers := strings.Split(string(output), "\n")
	for index, dockerName := range dockers {
		dockerName = strings.TrimSpace(dockerName)
		dockerName = strings.Trim(dockerName, "'")
		dockers[index] = dockerName
	}
	return dockers
}

func Start(name string) bool {
	cmd := exec.Command("docker", "start", name)
	output, err := cmd.CombinedOutput()
	fmt.Println(cmd.String())
	if err != nil {
		return false
	}
	return strings.Contains(string(output), name)
}

func StopAll() bool {
	dockers := List()
	for _, dockerName := range dockers {
		if strings.Contains(dockerName, "zentao") {
			Stop(dockerName)
		}
	}
	return true
}

func Stop(name string) bool {
	cmd := exec.Command("docker", "kill", name)
	fmt.Println(cmd.String())
	_, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return true
}

func InitZentao(version string) (err error) {
	versionNumber := strings.ReplaceAll(version, ".", "")
	containerName := "zentao" + versionNumber
	isExist := IsExistContainer(containerName)
	apath, _ := os.Getwd()
	codeDir := apath + "/docker/www/zentao" + versionNumber
	if runtime.GOOS == "windows" {
		codeDir = apath + `\docker\www\zentao` + versionNumber
	}
	if isExist {
		if !IsRuning(containerName) {
			StopAll()
			Start(containerName)
			time.Sleep(time.Second * 20)
		}
	} else {
		StopAll()
		err = Run(version, codeDir)
		if err != nil {
			return
		}
		time.Sleep(time.Second * 20)
	}
	err = uiTest.InitZentaoData(version, codeDir)
	return
}
