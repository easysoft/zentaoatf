import {app, BrowserWindow, Menu, shell} from 'electron';

import {IS_MAC_OSX, setEntryPath} from './utils/env';
import {logInfo, logErr} from './utils/log';
import {getUIServerUrl, startZtfServer, killZtfServer} from "./service";

export default class ZtfApp {
    constructor(entryPath) {
        this.startZtfServer()

        this._windows = new Map();
        setEntryPath(entryPath);
        this.bindElectronEvents();
        logInfo(`>> ZtfApp: created, entry path is "${entryPath}".`);
    }

    async startZtfServer() {
        try {
            const ztfServerUrl = await startZtfServer();
            logInfo(`>> ZTF Server started successfully: ${ztfServerUrl}`);
        } catch (error) {
            logErr('>> Start ztf server failed: ' + error);
            process.exit(1);
            return
        }
    }

    showAndFocus() {
        logInfo(`>> AppWindow[${this.name}]: show and focus`);

        const {browserWindow} = this;
        if (browserWindow.isMinimized()) {
            browserWindow.restore();
        } else {
            browserWindow.setOpacity(1);
            browserWindow.show();
        }
        browserWindow.focus();
    }

    createWindow() {
        process.env['ELECTRON_DISABLE_SECURITY_WARNINGS'] = 'true';

        require('@electron/remote/main').initialize()

        const mainWindow = new BrowserWindow({
            show: false,
            webPreferences: {nodeIntegration: true, contextIsolation: false}
        })
        require("@electron/remote/main").enable(mainWindow.webContents)

        mainWindow.maximize()
        mainWindow.show()

        this._windows.set('main', mainWindow);

        getUIServerUrl().then((url) => {
            mainWindow.loadURL(url);
            mainWindow.webContents.openDevTools({mode: 'bottom'});
        })
    };

    openOrCreateWindow() {
        const mainWin = this._windows.get('main');
        if (mainWin) {
            this.showAndFocus(mainWin)
        } else {
            this.createWindow();
        }
    }

    showAndFocus(mainWin) {
        if (mainWin.isMinimized()) {
            mainWin.restore();
        } else {
            mainWin.setOpacity(1);
            mainWin.show();
        }
        mainWin.focus();
    }

    ready() {
        logInfo('>> ZtfApp: ready.');
        this.openOrCreateWindow();
        this.buildAppMenu();
    }

    quit() {
        killZtfServer();
    }

    bindElectronEvents() {
        app.on('window-all-closed', () => {
            logInfo(`>> Event window-all-closed`)
            app.quit();
        });

        app.on('quit', () => {
            logInfo(`>> Event quit`)
            this.quit();
        });

        app.on('activate', () => {
            logInfo('>> ElectronApp: activate');

            // 在 OS X 系统上，可能存在所有应用窗口关闭了，但是程序还没关闭，此时如果收到激活应用请求需要
            // 重新打开应用窗口并创建应用菜单
            this.openOrCreateWindow();
            this.buildAppMenu();
        });
    }

    get windows() {
        return this._windows;
    }

    buildAppMenu() {
        if (!IS_MAC_OSX) {
            return;
        }

        const template = [
            {
                label: 'ZTF',
                submenu: [
                    {
                        label: '关于',
                        selector: 'orderFrontStandardAboutPanel:'
                    }, {
                        label: '退出',
                        accelerator: 'Command+Q',
                        click: () => {
                            this.quit();
                        }
                    }
                ]
            },
            {
                label: '编辑',
                submenu: [{
                    label: '撤销',
                    accelerator: 'Command+Z',
                    selector: 'undo:'
                }, {
                    label: '重做',
                    accelerator: 'Shift+Command+Z',
                    selector: 'redo:'
                }, {
                    type: 'separator'
                }, {
                    label: '剪切',
                    accelerator: 'Command+X',
                    selector: 'cut:'
                }, {
                    label: '复制',
                    accelerator: 'Command+C',
                    selector: 'copy:'
                }, {
                    label: '黏贴',
                    accelerator: 'Command+V',
                    selector: 'paste:'
                }, {
                    label: '选择所有',
                    accelerator: 'Command+A',
                    selector: 'selectAll:'
                }]
            },
            {
                label: '查看',
                submenu:  [
                    {
                        label: '切换全屏',
                        accelerator: 'Ctrl+Command+F',
                        click: () => {
                            const mainWin = this._windows.get('main');
                            mainWin.browserWindow.setFullScreen(!mainWin.browserWindow.isFullScreen());
                        }
                    }
                ]
            },
            {
                label: '窗口',
                submenu: [
                    {
                        label: '最小化',
                        accelerator: 'Command+M',
                        selector: 'performMiniaturize:'
                    },
                    {
                        label: '关闭',
                        accelerator: 'Command+W',
                        selector: 'performClose:'
                    },
                    {
                        type: 'separator'
                    },
                    {
                        label: '全部置于顶层',
                        selector: 'arrangeInFront:'
                    }
                ]
            },
            {
                label: '帮助',
                submenu: [{
                    label: '网站',
                    click: () => {
                        shell.openExternal('http://ztf.im');
                    }
                }]
            }];

        const menu = Menu.buildFromTemplate(template);
        Menu.setApplicationMenu(menu);

        logInfo('>> ZtfApp: build application menu.');
    }
}