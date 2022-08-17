package execHelper

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	configHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/config"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	commonUtils "github.com/easysoft/zentaoatf/pkg/lib/common"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func RunFile(filePath, workspacePath string, conf commDomain.WorkspaceConf,
	ch chan int, wsMsg *websocket.Message) (
	stdOutput string, errOutput string) {

	key := stringUtils.Md5(filePath)

	lang := langHelper.GetLangByFile(filePath)

	var cmd *exec.Cmd
	_, _, _, _, timeout := scriptHelper.GetCaseInfo(filePath)
	if timeout == 0 {
		timeout = 86400 * 7
	}
	ctxt, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	if commonUtils.IsWin() {
		scriptInterpreter := ""
		if strings.ToLower(lang) != "bat" {
			scriptInterpreter = configHelper.GetFieldVal(conf, stringUtils.UcFirst(lang))
		}
		if scriptInterpreter != "" {
			if strings.Index(strings.ToLower(scriptInterpreter), "autoit") > -1 {
				cmd = exec.Command("cmd", "/C", scriptInterpreter, filePath, "|", "more")
			} else {
				cmd = exec.Command("cmd", "/C", scriptInterpreter, filePath)
			}
		} else if strings.ToLower(lang) == "bat" {
			cmd = exec.Command("cmd", "/C", filePath)
		} else {
			msg := i118Utils.I118Prt.Sprintf("no_interpreter_for_run", lang, filePath)
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(msg, "", iris.Map{"key": key}, wsMsg)
			}
			logUtils.ExecConsolef(-1, msg)
			logUtils.ExecFilef(msg)
		}
	} else {
		err := os.Chmod(filePath, 0777)
		if err != nil {
			msg := i118Utils.I118Prt.Sprintf("exec_cmd_fail", filePath, err.Error())
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(msg, "", iris.Map{"key": key}, wsMsg)
			}
			logUtils.ExecConsolef(-1, msg)
			logUtils.ExecFilef(msg)
		}

		//filePath = "\"" + filePath + "\""
		scriptInterpreter := configHelper.GetFieldVal(conf, stringUtils.UcFirst(lang))

		if scriptInterpreter != "" {
			msg := fmt.Sprintf("use interpreter %s", scriptInterpreter)

			if commConsts.ExecFrom != commConsts.FromCmd {
				//websocketHelper.SendOutputMsg(msg, "", iris.Map{"key": key}, wsMsg)
				logUtils.ExecConsolef(-1, msg)
			}
			//logUtils.ExecFilef(msg)

			cmd = exec.CommandContext(ctxt, scriptInterpreter, filePath)
		} else {
			cmd = exec.CommandContext(ctxt, "/bin/bash", "-c", filePath)
		}
	}

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(msgStr, "", iris.Map{"key": key}, wsMsg)
			logUtils.ExecConsolef(color.FgRed, msgStr)
		}

		logUtils.ExecFilef(msgStr)

		return "", msgStr
	}

	cmd.Dir = workspacePath

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(err1.Error(), "", iris.Map{"key": key}, wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, err1.Error())
		logUtils.ExecFilef(err1.Error())

		return "", err1.Error()
	} else if err2 != nil {
		if commConsts.ExecFrom != commConsts.FromCmd {
			websocketHelper.SendOutputMsg(err2.Error(), "", iris.Map{"key": key}, wsMsg)
		}
		logUtils.ExecConsolef(color.FgRed, err2.Error())
		logUtils.ExecFilef(err2.Error())

		return "", err2.Error()
	}

	cmd.Start()

	isTerminal := false
	reader1 := bufio.NewReader(stdout)
	stdOutputArr := make([]string, 0)

	for {
		line, err2 := reader1.ReadString('\n')
		if line != "" {
			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendOutputMsg(line, "", iris.Map{"key": key}, wsMsg)
				logUtils.ExecConsole(-1, line)
			}

			logUtils.ExecFile(line)

			isTerminal = true
		}

		if err2 == io.EOF {
			break
		}
		if err2 != nil {
			logUtils.ExecConsole(color.FgRed, err2.Error())
			logUtils.ExecFile(err2.Error())
			break
		}

		stdOutputArr = append(stdOutputArr, line)

		select {
		case <-ch:
			msg := i118Utils.Sprintf("exit_exec_curr")

			if commConsts.ExecFrom != commConsts.FromCmd {
				websocketHelper.SendExecMsg(msg, "", commConsts.Run,
					iris.Map{"key": key, "status": "end"}, wsMsg)
			}

			logUtils.ExecConsolef(color.FgCyan, msg)
			logUtils.ExecFilef(msg)

			goto ExitCurrCase
		default:
		}
	}

ExitCurrCase:
	errOutputArr := make([]string, 0)
	if !isTerminal {
		reader2 := bufio.NewReader(stderr)

		for {
			line, err2 := reader2.ReadString('\n')
			if err2 != nil || io.EOF == err2 {
				break
			}
			errOutputArr = append(errOutputArr, line)
		}
	}
	if ctxt.Err() != nil {
		errOutputArr = append(errOutputArr, i118Utils.Sprintf("exec_cmd_timeout"))

	}
	stdOutput = strings.Join(stdOutputArr, "")
	errOutput = strings.Join(errOutputArr, "")
	cmd.Wait()
	return
}
