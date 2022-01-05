import path from 'path';
import {spawn} from 'child_process';
import {app} from 'electron';
import express from 'express';

let _ztfServerProcess;

export function startZtfServer() {
    if (process.env.SKIP_SERVER) {
        return Promise.resolve();
    }
    if (_ztfServerProcess) {
        return Promise.resolve(_ztfServerProcess);
    }

    let {SERVER_EXE_PATH: serverExePath} = process.env;
    if (serverExePath) {
        if (!path.isAbsolute(serverExePath)) {
            serverExePath = path.resolve(app.getAppPath(), serverExePath);
        }
        return new Promise((resolve) => {
            _ztfServerProcess = spawn(serverExePath);
            _ztfServerProcess.on('close', () => {
                _ztfServerProcess = null;
            });
            resolve(_ztfServerProcess);
        });
    }

    return new Promise((resolve, reject) => {
        const cmd = spawn('go', ['run', 'main.go', '-p', '8085'], {
            cwd: path.resolve(app.getAppPath(), '../cmd/server'),
            shell: true,
        });
        cmd.on('close', () => {
            _ztfServerProcess = null;
        });
        cmd.stdout.on('data', data => {
            const dataString = String(data);
            const lines = dataString.split('\n');
            for (let i = 0; i < lines.length; i++) {
                const line = lines[i];
                if (line.includes('Now listening on: http')) {
                    resolve(line.split('Now listening on:')[1].trim());
                    break;
                } else if (line.startsWith('[ERRO]')) {
                    reject(new Error(`Start ztf server failed with error: ${line.substring('[ERRO]'.length)}`));
                    break;
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
    if (uiServerUrl) {
        if (/^https?:\/\//.test(uiServerUrl)) {
            return Promise.resolve(uiServerUrl);
        }
        return new Promise((resolve, reject) => {
            if (!path.isAbsolute(uiServerUrl)) {
                uiServerUrl = path.resolve(app.getAppPath(), uiServerUrl);
            }

            console.log(`>> Starting UI serer at ${uiServerUrl}`);

            const uiServer = express();
            uiServer.use(express.static(uiServerUrl));
            const server = uiServer.listen(process.env.UI_SERVER_PORT || 8000, serverError => {
                if (serverError) {
                    console.error('>>> Start ui server failed with error', serverError);
                    _uiServerApp = null;
                    reject(serverError);
                } else {
                    const address = server.address();
                    console.log(`>> UI server started successfully on http://localhost:${address.port}.`);
                    resolve(`http://localhost:${address.port}`);
                }
            });
            server.on('close', () => {
                _uiServerApp = null;
            });
            _uiServerApp = uiServer;
        })
    }

    return new Promise((resolve, reject) => {
        console.log('>> Starting UI development serve...');

        let resolved = false;
        const cmd = spawn('npm', ['run', 'serve'], {
            cwd: path.resolve(app.getAppPath(), '../ui'),
            shell: true,
        });
        cmd.on('close', () => {
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
                if (line.includes('App running at:')) {
                    const nextLine = lines[i + 1] || lines[i + 2];
                    const url = nextLine.split('Local:   ')[1];
                    if (url) {
                        resolved = true;
                        resolve(url);
                    }
                    break;
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
