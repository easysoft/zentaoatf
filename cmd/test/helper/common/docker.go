package commonTestHelper

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"time"

	constTestHelper "github.com/easysoft/zentaoatf/cmd/test/helper/conf"
	uiTest "github.com/easysoft/zentaoatf/cmd/test/helper/zentao/ui"
)

func Run(version string) (err error) {
	versionNumber := strings.ReplaceAll(version, ".", "_")

	cmd := exec.Command("docker", "run", "--name", "zentao"+versionNumber, "-p",
		fmt.Sprintf("%d:80", constTestHelper.ZentaoPort), "-e", "MYSQL_INTERNAL=true", "-d", "hub.zentao.net/app/zentao:"+version)
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

func IsRunning(name string) bool {
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
	versionNumber := strings.ReplaceAll(version, ".", "_")
	containerName := "zentao" + versionNumber
	isExist := IsExistContainer(containerName)

	if isExist {
		if !IsRunning(containerName) {
			StopAll()
			Start(containerName)

			waitZentaoAccessed()
		}
	} else {
		StopAll()
		err = Run(version)
		if err != nil {
			return
		}

		waitZentaoAccessed()
	}

	err = uiTest.InitZentaoData(version)
	return
}

func InitZentaoData() (err error) {
	err = uiTest.InitZentaoData("latest")
	return
}

func waitZentaoAccessed() {
	isTimeout := false
	time.AfterFunc(80*time.Second, func() {
		isTimeout = true
	})

	for {
		status := uiTest.GetStatus(constTestHelper.ZentaoSiteUrl)
		if isTimeout || status {
			return
		}
		time.Sleep(3 * time.Second)
		fmt.Println("waiting zentao ...")
	}
}
