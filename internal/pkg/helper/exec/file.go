package execHelper

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"strings"
	"time"

	scriptHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/script"

	commConsts "github.com/easysoft/zentaoatf/internal/pkg/consts"
	commDomain "github.com/easysoft/zentaoatf/internal/pkg/domain"
	langHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/lang"
	websocketHelper "github.com/easysoft/zentaoatf/internal/pkg/helper/websocket"
	i118Utils "github.com/easysoft/zentaoatf/pkg/lib/i118"
	logUtils "github.com/easysoft/zentaoatf/pkg/lib/log"
	stringUtils "github.com/easysoft/zentaoatf/pkg/lib/string"
	"github.com/fatih/color"
	"github.com/gofrs/uuid"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/websocket"
)

func RunFile(filePath, workspacePath string, conf commDomain.WorkspaceConf,
	ch chan int, wsMsg *websocket.Message, idx int) (
	stdOutput string, errOutput string) {

	key := stringUtils.Md5(filePath)

	lang := langHelper.GetLangByFile(filePath)

	uuidString := uuid.Must(uuid.NewV4()).String()
	_, _, _, _, timeout := scriptHelper.GetCaseInfo(filePath)
	if timeout == 0 {
		timeout = 86400 * 7
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	cmd := getCommand(filePath, lang, uuidString, conf, ctx, wsMsg)

	if cmd == nil {
		msgStr := i118Utils.Sprintf("cmd_empty")
		if commConsts.ExecFrom == commConsts.FromClient {
			websocketHelper.SendOutputMsg(msgStr, "", iris.Map{"key": key}, wsMsg)
			logUtils.ExecConsolef(color.FgRed, msgStr)
		}

		logUtils.ExecFilef(msgStr)

		return "", msgStr
	} else {
		cmd.Dir = workspacePath
	}

	cmd.Dir = workspacePath
	if commConsts.BatchCount > 1 {
		cmd.Env = append(cmd.Env, "ZTF_POOL_ID="+strconv.Itoa(idx+1))
	}
	cmd.Env = append(cmd.Env, "ZTF_REPORT_DIR="+commConsts.ExecLogDir)

	stdout, err1 := cmd.StdoutPipe()
	stderr, err2 := cmd.StderrPipe()

	if err1 != nil {
		return PrintErrMsg(key, err1, wsMsg)
	} else if err2 != nil {
		return PrintErrMsg(key, err2, wsMsg)
	}

	cmd.Start()
	isTerminal := false

	go func() {
		time.AfterFunc(time.Second*time.Duration(timeout), func() {
			KillProcessByUUID(uuidString)
			isTerminal = true
			stdout.Close()
			stderr.Close()
		})
		if ch != nil {
			for {
				if isTerminal {
					break
				}
				select {
				case _, ok := <-ch:
					KillProcessByUUID(uuidString)
					stdout.Close()
					stderr.Close()
					SetRunning(false)
					if ok {
						close(ch)
					}
					return
				default:
					time.Sleep(time.Millisecond * 100)
				}
			}
		}

	}()
	reader1 := bufio.NewReader(stdout)
	stdOutputArr := make([]string, 0)

	for {
		line, err2 := reader1.ReadString('\n')
		if line != "" {
			if commConsts.ExecFrom == commConsts.FromClient {
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

			if commConsts.ExecFrom == commConsts.FromClient {
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
	if ctx.Err() != nil {
		errOutputArr = append(errOutputArr, i118Utils.Sprintf("exec_cmd_timeout"))

	}
	stdOutput = strings.Join(stdOutputArr, "")
	errOutput = strings.Join(errOutputArr, "")
	cmd.Wait()
	return
}
