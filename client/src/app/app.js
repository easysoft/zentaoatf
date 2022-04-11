import {app, BrowserWindow, Menu, shell} from 'electron';
import main from "@electron/remote/main";

import {DEBUG} from './utils/consts';
import Config, {updateConfig} from './utils/config';
import {IS_MAC_OSX} from './utils/env';
import Lang, {initLang} from './core/lang';
import {logInfo, logErr} from './utils/log';
import {startUIService} from "./core/ui";
import {startZtfServer, killZtfServer} from "./core/ztf";

export default class ZtfApp {
    constructor() {
        app.name = Lang.string('app.title', Config.pkg.displayName);

        this._windows = new Map();

        startZtfServer().then((ztfServerUrl)=> {
            if (ztfServerUrl) logInfo(`>> ztf server started successfully on : ${ztfServerUrl}`);
            this.bindElectronEvents();
        }).catch((err) => {
            logErr('>> ztf server started failed, err: ' + err);
            process.exit(1);
            return;
        })
    }

    showAndFocus() {
        logInfo(`>> ztf app: AppWindow[${this.name}]: show and focus`);

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

        const remoteMain = require('@electron/remote/main')
        remoteMain.initialize()

        const mainWin = new BrowserWindow({
            show: false,
            webPreferences: {nodeIntegration: true, contextIsolation: false}
        })
        remoteMain.enable(mainWin.webContents)

        mainWin.maximize()
        mainWin.show()

        this._windows.set('main', mainWin);

        const url = await startUIService()
        await mainWin.loadURL(url);

        if (DEBUG) {
            mainWin.webContents.openDevTools({mode: 'bottom'});
        }
    };

    async openOrCreateWindow() {
        const mainWin = this._windows.get('main');
        if (mainWin) {
            this.showAndFocus(mainWin)
        } else {
            await this.createWindow();
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
        logInfo('>> ztf app ready.');

        initLang()
        this.buildAppMenu();
        this.openOrCreateWindow()
        this.setAboutPanel();
    }

    quit() {
        killZtfServer();
    }

    bindElectronEvents() {
        app.on('window-all-closed', () => {
            logInfo(`>> event: window-all-closed`)
            app.quit();
        });

        app.on('quit', () => {
            logInfo(`>> event: quit`)
            this.quit();
        });

        app.on('activate', () => {
            logInfo('>> event: activate');

            this.buildAppMenu();

            // 在 OS X 系统上，可能存在所有应用窗口关闭了，但是程序还没关闭，此时如果收到激活应用请求需要
            // 重新打开应用窗口并创建应用菜单
            this.openOrCreateWindow().then(() => {

            })
        });
    }

    get windows() {
        return this._windows;
    }

    buildAppMenu() {
        logInfo('>> ztf app: build application menu.');

        if (!IS_MAC_OSX) {
            return;
        }

        const template = [
            {
                label: Lang.string('app.title', Config.pkg.displayName),
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
                            mainWin.setFullScreen(!mainWin.isFullScreen());
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

    setAboutPanel() {
        if (!app.setAboutPanelOptions) {
            return;
        }
        app.setAboutPanelOptions({
            applicationName: Lang.string(Config.pkg.name) || Config.pkg.displayName,
            applicationVersion: Config.pkg.version,
            copyright: `${Config.pkg.copyright} ${Config.pkg.company}`,
            credits: `Licence: ${Config.pkg.license}`,
            version: `${Config.pkg.buildTime ? `build at ${new Date(Config.pkg.buildTime).toLocaleString()}` : ''}${DEBUG ? '[debug]' : ''}`
        });
    }
}