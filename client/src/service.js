import path from 'path';
import cp, {exec, spawn} from 'child_process';
import os from 'os';
import {app} from 'electron';
import express from 'express';
import {logInfo, logErr} from './log';

const psTree = require('ps-tree');
const DEBUG = process.env.NODE_ENV === 'development';
const isWin = /^win/.test(process.platform);
const isMac = /^darwin/.test(process.platform);
const uuid = 'ZTF@1CF17A46-B136-4AEB-96B4-F21C8200EF5A~'
const port = 8085
let _ztfServerProcess;
let _ztfSubProcessIds = [];

export function checkZtfPort() {
    let cmd = ''
    if (!isWin) {
        cmd = 'lsof -i:${port} | grep ${port}'
    } else {
        cmd = 'netstat -aon | findstr ${port}'
    }

    const cp = require('child_process');
    const stdout = cp.execSync(cmd).toString().trim()
    console.log('exec ${cmd}, stdout: ${stdout}');

    if (stdout.indexOf(port + '') > -1) {
        if (stdout.indexOf(uuid) < 0) {
            const msg = 'Port ${port} is used by another process. exit.'
            console.log(msg);
            logErr(msg);
            return false
        }

        const msg = 'Port ${port} is used by ztf process. kill previous one.'
        console.log(msg);
        logInfo(msg);

        killZtfServer()
    } else {
        return true
    }
}

export function startZtfServer() {
    if (process.env.SKIP_SERVER) {
        logInfo(`>> Skip to start ZTF Server by env "SKIP_SERVER=${process.env.SKIP_SERVER}".`);
        return Promise.resolve();
    }
    if (_ztfServerProcess) {
        return Promise.resolve(_ztfServerProcess);
    }

    let {SERVER_EXE_PATH: serverExePath} = process.env;
    if (!serverExePath && !DEBUG) {
        const platform = os.platform(); // 'darwin', 'linux', 'win32'
        const exePath = `bin/${platform}/ztf${platform === 'win32' ? '.exe' : ''}`;
        serverExePath = path.join(process.resourcesPath, exePath);
    }
    if (serverExePath) {
        if (!path.isAbsolute(serverExePath)) {
            serverExePath = path.resolve(app.getAppPath(), serverExePath);
        }
        return new Promise((resolve, reject) => {
            const cwd = process.env.SERVER_CWD_PATH || path.dirname(serverExePath);
            logInfo(`>> Starting ZTF Server from exe path with command "${serverExePath} -P 8085" in "${cwd}"...`);
            const cmd = spawn(serverExePath, ['-p', '8085', "-uuid", uuid], {
                cwd,
                shell: true,
            });
            cmd.on('close', (code) => {
                logInfo(`>> ZTF server closed with code ${code}`);
                _ztfServerProcess = null;
                cmd.kill()
            });
            cmd.stdout.on('data', data => {
                const dataString = String(data);
                const lines = dataString.split('\n');
                for (let i = 0; i < lines.length; i++) {
                    const line = lines[i];
                    if (DEBUG) {
                        logInfo('\t' + line);
                    }
                    if (line.includes('Now listening on: http')) {
                        resolve(line.split('Now listening on:')[1].trim());
                        if (!DEBUG) {
                            break;
                        }
                    } else if (line.includes('启动HTTP服务于')) {
                        resolve(line.split(/启动HTTP服务于|，/)[1].trim());
                        if (!DEBUG) {
                            break;
                        }
                    } else if (line.startsWith('[ERRO]')) {
                        reject(new Error(`Start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
                        if (!DEBUG) {
                            break;
                        }
                    }
                }
            });
            cmd.on('error', spawnError => {
                console.error('>>> Start ztf server failed with error', spawnError);
                reject(spawnError)
            });
            _ztfServerProcess = cmd;
            logInfo(`>> _ztfServerProcess = ${_ztfServerProcess.pid}`)

            psTree(_ztfServerProcess.pid, function (err, children) {
                _ztfSubProcessIds = [_ztfServerProcess.pid].concat(
                    children.map(function (p) {
                        return p.PID;
                    })
                );
                logInfo(`>> _ztfSubProcessIds = ${_ztfSubProcessIds}`)
            });
        });
    }

    return new Promise((resolve, reject) => {
        const cwd = process.env.SERVER_CWD_PATH || path.resolve(app.getAppPath(), '../');
        logInfo(`>> Starting ZTF development server from source with command "go run cmd/server/main.go -P 8085" in "${cwd}"`);
        const cmd = spawn('go', ['run', 'main.go', '-P', '8085'], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            logInfo(`>> ZTF server closed with code ${code}`);
            _ztfServerProcess = null;
        });
        cmd.stdout.on('data', data => {
            const dataString = String(data);
            const lines = dataString.split('\n');
            for (let i = 0; i < lines.length; i++) {
                const line = lines[i];
                if (DEBUG) {
                    logInfo('\t' + line);
                }
                if (line.includes('Now listening on: http')) {
                    resolve(line.split('Now listening on:')[1].trim());
                    if (!DEBUG) {
                        break;
                    }
                } else if (line.startsWith('[ERRO]')) {
                    reject(new Error(`Start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
                    if (!DEBUG) {
                        break;
                    }
                }
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> Start ztf server failed with error', spawnError);
            reject(spawnError)
        });
        _ztfServerProcess = cmd;
    });
}

let _uiServerApp;

export function getUIServerUrl() {
    if (_uiServerApp) {
        return Promise.resolve();
    }

    let {UI_SERVER_URL: uiServerUrl} = process.env;
    if (!uiServerUrl && !DEBUG) {
        uiServerUrl = path.resolve(process.resourcesPath, 'ui');
    }

    if (uiServerUrl) {
        if (/^https?:\/\//.test(uiServerUrl)) {
            return Promise.resolve(uiServerUrl);
        }
        return new Promise((resolve, reject) => {
            if (!path.isAbsolute(uiServerUrl)) {
                uiServerUrl = path.resolve(app.getAppPath(), uiServerUrl);
            }

            const port = process.env.UI_SERVER_PORT || 8000;
            logInfo(`>> Starting UI serer at ${uiServerUrl} with port ${port}`);

            const uiServer = express();
            uiServer.use(express.static(uiServerUrl));
            const server = uiServer.listen(port, serverError => {
                if (serverError) {
                    console.error('>>> Start ui server failed with error', serverError);
                    _uiServerApp = null;
                    reject(serverError);
                } else {
                    logInfo(`>> UI server started successfully on http://localhost:${port}.`);
                    resolve(`http://localhost:${port}`);
                }
            });
            server.on('close', () => {
                _uiServerApp = null;
            });
            _uiServerApp = uiServer;
        })
    }

    return new Promise((resolve, reject) => {
        const cwd = path.resolve(app.getAppPath(), '../ui');
        logInfo(`>> Starting UI development server with command "npm run serve" in "${cwd}"...`);

        let resolved = false;
        const cmd = spawn('npm', ['run', 'serve'], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            logInfo(`>> ZTF server closed with code ${code}`);
            _uiServerApp = null;
        });
        cmd.stdout.on('data', data => {
            if (resolved) {
                return;
            }
            const dataString = String(data);
            const lines = dataString.split('\n');
            for (let i = 0; i < lines.length; i++) {
                const line = lines[i];
                if (DEBUG) {
                    logInfo('\t' + line);
                }
                if (line.includes('App running at:')) {
                    const nextLine = lines[i + 1] || lines[i + 2];
                    if (DEBUG) {
                        logInfo('\t' + nextLine);
                    }
                    if (!nextLine) {
                        console.error('\t' + `Cannot grabing running address after line "${line}".`);
                        throw new Error(`Cannot grabing running address after line "${line}".`);
                    }
                    const url = nextLine.split('Local:   ')[1];
                    if (url) {
                        resolved = true;
                        resolve(url);
                    }
                    if (!DEBUG) {
                        break;
                    }
                }
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> Get ui server url failed with error', spawnError);
            reject(spawnError)
        });
        _uiServerApp = cmd;
    });
}

export function killZtfServer() {
    let cmd = ''
    if (!isWin) {
        logInfo(`>> no windows`);

        cmd = `ps -ef | grep ${uuid} | grep -v "\-\-%s" | grep -v "grep" | awk '{print $2}' | xargs kill -9`
        logInfo(`kill cmd : ${cmd}`);
        const cp = require('child_process');
        cp.exec(cmd, function (error, stdout, stderr) {
            logInfo(`stdout: ${stdout}; stderr: ${stderr}; error: ${error}`);
        });
    } else {
        logInfo(`>> is windows`);
			
        const cmd = 'WMIC path win32_process  where "Commandline like \'%%' + uuid + '%%\'" get Processid,Caption';
        let msg = `list process cmd : ${cmd}`
        console.log(msg);
        logInfo(msg);

        const cp = require('child_process');
        const stdout = cp.execSync(cmd).toString().trim()
        msg = `exec ${cmd}, stdout: ${stdout}`
        console.log(msg);
        logInfo(msg)

        const lines = stdout.split('\n')
        lines.forEach(function(line){
            line = line.trim()
            console.log(`<${line}>`)
            logInfo(`<${line}>`)
            const columns = line.split(/\s/)

            if (columns.length > 2) {
                const pid = columns[2].trim()
                console.log(`pid=${pid}`);
                logInfo(`pid=${pid}`)

                if (pid && parseInt(pid, 10) === pid) {
                    const killCmd = `taskkill /F /pid ${pid}`

                    const cp = require('child_process');
                    const stdout = cp.execSync(killCmd).toString().trim()
                    msg = `exec ${cmd}, stdout: ${stdout}`
                    console.log(msg);
                    logInfo(msg)
                }
            }
        });
    }

    // if (isWin) {
    //     logInfo(`>> isWin`);
    //     const cp = require('child_process');
    //     cp.exec('taskkill /PID ' + _ztfServerProcess.pid + ' /T /F',
    //         function (error, stdout, stderr) {
    //             // logInfo('stdout: ' + stdout + '; stderr: ' + stderr + '; error: ' + error + '.');
    //         });
    // } else if (isMac) {
    //     logInfo(`>> isMac`);
    //     if (_ztfServerProcess) _ztfServerProcess.kill();
    // } else {
    //     logInfo(`>> isLinux`);
    //     _ztfSubProcessIds.forEach(function (tpid) {
    //         logInfo(`>> kill ${tpid}`)
    //         try {
    //             process.kill(tpid, 'SIGKILL')
    //         } catch (ex) {
    //             logErr(`>> ` + ex)
    //         }
    //     });
    // }
}


