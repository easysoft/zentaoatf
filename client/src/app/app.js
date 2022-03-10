import {app, BrowserWindow, Menu, shell} from 'electron';

import {DEBUG} from './utils/consts';
import {IS_MAC_OSX, setEntryPath} from './utils/env';
import Lang, {initLang} from './core/lang';

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
            logInfo(`>> ZtfServer started successfully: ${ztfServerUrl}`);
        } catch (error) {
            logErr('>> ZtfServer started failed: ' + error);
            process.exit(1);
            return;
        }
    }

    showAndFocus() {
        logInfo(`>> ZtfApp: AppWindow[${this.name}]: show and focus`);

        const {browserWindow} = this;
        if (browserWindow.isMinimized()) {
            browserWindow.restore();
        } else {
            browserWindow.setOpacity(1);
            browserWindow.show();
        }
        browserWindow.focus();
    }

    async createWindow() {
        process.env['ELECTRON_DISABLE_SECURITY_WARNINGS'] = 'true';

        require('@electron/remote/main').initialize()

        const mainWin = new BrowserWindow({
            show: false,
            webPreferences: {nodeIntegration: true, contextIsolation: false}
        })
        require("@electron/remote/main").enable(mainWin.webContents)

        mainWin.maximize()
        mainWin.show()

        this._windows.set('main', mainWin.l);

        const url = await getUIServerUrl()
        await mainWin.loadURL(url);
        mainWin.webContents.openDevTools({mode: 'bottom'});
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

    async ready() {
        logInfo('>> ZtfApp: ready.');

        initLang()
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
        logInfo('>> ZtfApp: build application menu.');

        if (!IS_MAC_OSX) {
            return;
        }

        const template = [
            {
                label: 'ZTF',
                submenu: [
                    {
                        label: Lang.string('app.about'),
                        selector: 'orderFrontStandardAboutPanel:'
                    }, {
                        label: Lang.string('app.exit'),
                        accelerator: 'Command+Q',
                        click: () => {
                            app.quit();
                        }
                    }
                ]
            },
            {
                label: Lang.string('app.edit'),
                submenu: [{
                    label: Lang.string('app.undo'),
                    accelerator: 'Command+Z',
                    selector: 'undo:'
                }, {
                    label: Lang.string('app.redo'),
                    accelerator: 'Shift+Command+Z',
                    selector: 'redo:'
                }, {
                    type: 'separator'
                }, {
                    label: Lang.string('app.cut'),
                    accelerator: 'Command+X',
                    selector: 'cut:'
                }, {
                    label: Lang.string('app.copy'),
                    accelerator: 'Command+C',
                    selector: 'copy:'
                }, {
                    label: Lang.string('app.paste'),
                    accelerator: 'Command+V',
                    selector: 'paste:'
                }, {
                    label: Lang.string('app.select_all'),
                    accelerator: 'Command+A',
                    selector: 'selectAll:'
                }]
            },
            {
                label: Lang.string('app.view'),
                submenu:  [
                    {
                        label: Lang.string('app.switch_to_full_screen'),
                        accelerator: 'Ctrl+Command+F',
                        click: () => {
                            const mainWin = this._windows.get('main');
                            mainWin.browserWindow.setFullScreen(!mainWin.browserWindow.isFullScreen());
                        }
                    }
                ]
            },
            {
                label: Lang.string('app.window'),
                submenu: [
                    {
                        label: Lang.string('app.minimize'),
                        accelerator: 'Command+M',
                        selector: 'performMiniaturize:'
                    },
                    {
                        label: Lang.string('app.close'),
                        accelerator: 'Command+W',
                        selector: 'performClose:'
                    },
                    {
                        type: 'separator'
                    },
                    {
                        label: Lang.string('app.bring_all_to_front'),
                        selector: 'arrangeInFront:'
                    }
                ]
            },
            {
                label: Lang.string('app.help'),
                submenu: [{
                    label: Lang.string('app.website'),
                    click: () => {
                        shell.openExternal('http://ztf.im');
                    }
                }]
            }];

        const menu = Menu.buildFromTemplate(template);
        Menu.setApplicationMenu(menu);
    }
}