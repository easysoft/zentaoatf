import {app} from 'electron';
import path from 'path';
import {spawn} from 'child_process';

let _ztfServerProcess;

export function startZtfServer() {
    if (_ztfServerProcess) {
        return Promise.resolve(_ztfServerProcess);
    }

    const {SERVER_EXE_PATH} = process.env;
    if (SERVER_EXE_PATH) {
        return new Promise((resolve) => {
            _ztfServerProcess = spawn(SERVER_EXE_PATH);
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
            if (dataString.includes('Now listening on: http')) {
                resolve(_ztfServerProcess);
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> Start ztf server failed with error', spawnError);
            reject(spawnError)
        });
        _ztfServerProcess = cmd;
    });
}

let _uiServerProcess;

export function getUIServerUrl() {
    if (_uiServerProcess) {
        return Promise.resolve(_uiServerProcess);
    }

    const {UI_SERVER_PATH} = process.env;
    if (UI_SERVER_PATH) {
        return new Promise((resolve) => {
            _uiServerProcess = spawn(UI_SERVER_PATH);
            _uiServerProcess.on('close', () => {
                _uiServerProcess = null;
            });
            resolve(_uiServerProcess);
        });
    }

    return new Promise((resolve, reject) => {
        let resolved = false;
        const cwd = path.resolve(app.getAppPath(), '../ui');
        console.log('Get UI server url with cwd: ', cwd);
        const cmd = spawn('npm', ['run', 'serve'], {
            cwd,
            shell: true,
        });
        cmd.on('close', () => {
            _uiServerProcess = null;
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
                    const nextLine = lines[i + 1];
                    const url = nextLine.split('Local:   ')[1];
                    resolved = true;
                    resolve(url);
                    break;
                }
            }
        });
        cmd.on('error', spawnError => {
            console.error('>>> Get ui server url failed with error', spawnError);
            reject(spawnError)
        });
        _uiServerProcess = cmd;
    });
}
