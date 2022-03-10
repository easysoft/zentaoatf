import path from 'path';
import {execSync, spawn} from 'child_process';
import os from 'os';
import {app} from 'electron';

import {IS_WINDOWS_OS} from "../utils/env";
import {DEBUG, portServer, uuid} from '../utils/consts';
import {logErr, logInfo} from '../utils/log';

let _ztfProcess;

export async function startZtfServer() {
    // await killPortProcess([portClient, portServer])
    if (process.env.SKIP_SERVER) {
        logInfo(`>> skip to start ztf Server by env "SKIP_SERVER=${process.env.SKIP_SERVER}".`);
        return Promise.resolve();
    }
    if (_ztfProcess) {
        return Promise.resolve(_ztfProcess);
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
            logInfo(`>> starting ztf server with command ` +
                `"${serverExePath} -p ${portServer} -uuid ${uuid}" in "${cwd}"...`);

            const cmd = spawn(serverExePath, ['-p', portServer, '-uuid', uuid], {
                cwd,
                shell: true,
            });

            _ztfProcess = cmd;
            logInfo(`>> ztf server process = ${_ztfProcess.pid}`)

            cmd.on('close', (code) => {
                logInfo(`>> ztf server closed with code ${code}`);
                _ztfProcess = null;
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
                        reject(new Error(`start ztf server failed, error: ${line.substring('[ERRO]'.length)}`));
                        if (!DEBUG) {
                            break;
                        }
                    }
                }
            });
            cmd.on('error', spawnError => {
                logErr('>>> start ztf server failed with error', spawnError);
                reject(spawnError)
            });
        });
    }

    return new Promise((resolve, reject) => {
        const cwd = process.env.SERVER_CWD_PATH || path.resolve(app.getAppPath(), '../');
        logInfo(`>> starting ztf development server from source with command "go run cmd/server/main.go -p ${portServer}" in "${cwd}"`);
        const cmd = spawn('go', ['run', 'main.go', '-p', portServer], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            logInfo(`>> ztf server closed with code ${code}`);
            _ztfProcess = null;
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
                    reject(new Error(`start ztf server failed, error: ${line.substring('[ERRO]'.length)}`));
                    if (!DEBUG) {
                        break;
                    }
                }
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> start ztf server failed with error', spawnError);
            reject(spawnError)
        });
        _ztfProcess = cmd;
    });
}

export function killZtfServer() {
    let cmd = ''
    if (!IS_WINDOWS_OS) {
        logInfo(`>> no windows`);

        cmd = `ps -ef | grep ${uuid} | grep -v "grep" | awk '{print $2}' | xargs kill -9`
        logInfo(`>> kill cmd: ${cmd}`);

        const cp = require('child_process');
        cp.exec(cmd, function (error, stdout, stderr) {
            logInfo(`>> kill result: stdout: ${stdout}; stderr: ${stderr}; error: ${error}`);
        });
    } else {
        const cmd = 'WMIC path win32_process  where "Commandline like \'%%' + uuid + '%%\'" get Processid,Caption';
        logInfo(`>> list process cmd: ${cmd}`);

        const stdout = execSync(cmd, {windowsHide: true}).toString().trim()
        logInfo(`>> list process result: exec ${cmd}, stdout: ${stdout}`)

        let pid = 0
        const lines = stdout.split('\n')
        lines.forEach(function(line){
            line = line.trim()
            console.log(`<${line}>`)
            logInfo(`<${line}>`)
            const cols = line.split(/\s/)

            if (line.indexOf('ztf') > -1 && cols.length > 3) {
                const col3 = cols[3].trim()
                console.log(`col3=${col3}`);
                logInfo(`col3=${col3}`)

                if (col3 && parseInt(col3, 10)) {
                    pid = parseInt(col3, 10)
                }
            }
        });

        if (pid && pid > 0) {
            const killCmd = `taskkill /F /pid ${pid}`
            logInfo(`>> taskkill cmd: exec ${killCmd}`)

            const out = execSync(`taskkill /F /pid ${pid}`, {windowsHide: true}).toString().trim()
            logInfo(`>> taskkill result: ${out}`)
        }
    }
}