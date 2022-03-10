import path from 'path';
import {execSync, spawn} from 'child_process';
import os from 'os';
import {app} from 'electron';
import express from 'express';

import {DEBUG, portClient, portServer, uuid} from './utils/consts';
import {logInfo, logDebug} from './utils/log';
import {IS_WINDOWS_OS} from "./utils/env";

let _ztfServerProcess;

export async function startZtfServer() {
    // await killPortProcess([portClient, portServer])
    if (process.env.SKIP_SERVER) {
        logInfo(`>> skip to start ztf Server by env "SKIP_SERVER=${process.env.SKIP_SERVER}".`);
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
            logInfo(`>> starting ztf server with command ` +
                `"${serverExePath} -p ${portServer} -uuid ${uuid}" in "${cwd}"...`);

            const cmd = spawn(serverExePath, ['-p', portServer, '-uuid', uuid], {
                cwd,
                shell: true,
            });

            _ztfServerProcess = cmd;
            logInfo(`>> ztf server process = ${_ztfServerProcess.pid}`)

            cmd.on('close', (code) => {
                logInfo(`>> ztf server closed with code ${code}`);
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
                    if (line.includes('now listening on: http')) {
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
                        reject(new Error(`start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
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
                if (line.includes('now listening on: http')) {
                    resolve(line.split('now listening on:')[1].trim());
                    if (!DEBUG) {
                        break;
                    }
                } else if (line.startsWith('[ERRO]')) {
                    reject(new Error(`start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
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

            const port = process.env.UI_SERVER_PORT || portClient;
            logInfo(`>> starting ui serer at ${uiServerUrl} with port ${port}`);

            const uiServer = express();
            uiServer.use(express.static(uiServerUrl));
            const server = uiServer.listen(port, serverError => {
                if (serverError) {
                    console.error('>>> start ui server failed with error', serverError);
                    _uiServerApp = null;
                    reject(serverError);
                } else {
                    logInfo(`>> ui server started successfully on http://localhost:${port}.`);
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
        logInfo(`>> starting ui development server with command "npm run serve" in "${cwd}"...`);

        let resolved = false;
        const cmd = spawn('npm', ['run', 'serve'], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            logInfo(`>> ztf server closed with code ${code}`);
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
            console.error('>>> get ui server url failed with error', spawnError);
            reject(spawnError)
        });
        _uiServerApp = cmd;
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


