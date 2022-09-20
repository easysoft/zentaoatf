package commonTest

import (
	"errors"
	"flag"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	uiTest "github.com/easysoft/zentaoatf/test/helper/zentao/ui"
)

var version = flag.String("version", "", "")

func Run(version string) (err error) {
	versionNumber := strings.ReplaceAll(version, ".", "")
	apath, _ := os.Getwd()
	codeDir := apath + "/docker/www"
	if runtime.GOOS == "windows" {
		codeDir = apath + `\docker\www`
	}

	_, err = os.Stat(codeDir)
	if os.IsExist(err) {
		os.RemoveAll(codeDir)
	}

	//docker run --name zentao -p 8081:80 --network=zentaonet -v D:\docker\www\zentaopms:/www/zentaopms -v D:\docker\mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d easysoft/zentao:12.3.3
	cmd := exec.Command("docker", "run", "--name", "zentao"+versionNumber, "-p", "8081:80", "-v", codeDir+":/www/zentaopms", "-e", "MYSQL_ROOT_PASSWORD=123456", "-d", "easysoft/zentao:"+version)
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

func Start(name string) bool {
	cmd := exec.Command("docker", "start", name)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}
	return strings.Contains(string(output), name)
}

func InitZentao() {
	versionNumber := strings.ReplaceAll(*version, ".", "")
	containerName := "zentao" + versionNumber
	isExist := IsExistContainer(containerName)
	if isExist {
		if !IsRuning(containerName) {
			Start(containerName)
			time.Sleep(time.Second * 30)
		}
	} else {
		Run(*version)
		time.Sleep(time.Second * 120)
	}
	uiTest.InitZentaoData()
}
