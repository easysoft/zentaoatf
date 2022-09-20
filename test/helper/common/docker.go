package commonTest

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	uiTest "github.com/easysoft/zentaoatf/test/helper/zentao/ui"
)

var version = flag.String("version", "", "")
var isRuning = false

func Run(version string) (err error) {
	versionNumber := strings.ReplaceAll(version, ".", "")
	apath, _ := os.Getwd()
	fmt.Println(apath)
	codeDir := apath + "/docker/www"
	if runtime.GOOS == "windows" {
		codeDir = apath + `\docker\www`
	}

	_, err = os.Stat(codeDir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(codeDir, os.ModePerm)
		if err != nil {
			return
		}
	}

	//docker run --name zentao -p 8081:80 --network=zentaonet -v D:\docker\www\zentaopms:/www/zentaopms -v D:\docker\mysql:/var/lib/mysql -e MYSQL_ROOT_PASSWORD=123456 -d easysoft/zentao:12.3.3
	cmd := exec.Command("docker", "run", "--name", "zentao"+versionNumber, "-p", "8081:80", "-v", codeDir+":/www/zentaopms", "-e", "MYSQL_ROOT_PASSWORD=123456", "-d", "easysoft/zentao:"+version)
	fmt.Println(cmd.String())
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(output))
	return err
}

func InitZentao() {
	if isRuning {
		return
	}
	flag.Parse()
	isRuning = true
	fmt.Println(*version)
	Run(*version)
	time.Sleep(time.Minute)
	uiTest.InitZentaoData()
}
