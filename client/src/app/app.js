import {app, BrowserWindow, ipcMain, Menu, shell, dialog, globalShortcut} from 'electron';

import {DEBUG, electronMsg, electronMsgReplay, minimumSizeHeight, minimumSizeWidth} from './utils/consts';
import {IS_MAC_OSX, IS_LINUX, IS_WINDOWS_OS} from './utils/env';
import {logInfo, logErr} from './utils/log';
import Config from './utils/config';
import Lang, {initLang} from './core/lang';
import {startUIService} from "./core/ui";
import {startZtfServer, killZtfServer} from "./core/ztf";

const cp = require('child_process');
const fs = require('fs');
const pth = require('path');

export class ZtfApp {
    constructor() {
        app.name = Lang.string('app.title', Config.pkg.displayName);

        this._windows = new Map();

        killZtfServer();

        startZtfServer().then((ztfServerUrl)=> {
            if (ztfServerUrl) logInfo(`>> ztf server started successfully on : ${ztfServerUrl}`);
            this.bindElectronEvents();
        }).catch((err) => {
            logErr('>> ztf server started failed, err: ' + err);
            process.exit(1);
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

        const mainWin = new BrowserWindow({
            show: false,
            frame: false,
            webPreferences: {
                nodeIntegration: true,
                contextIsolation: false,
            }
        })

        mainWin.setSize(minimumSizeWidth, minimumSizeHeight)
        mainWin.setMovable(true)
        mainWin.maximize()
        mainWin.show()

        this._windows.set('main', mainWin);

        const url = await startUIService()
        await mainWin.loadURL(url);

        ipcMain.on(electronMsg, (event, arg) => {
            logInfo('msg from renderer: ' + arg)

            switch (arg) {
                case 'selectDir':
                    this.showFolderSelection(event)
                    break;
                case 'selectFile':
                    this.showFileSelection(event)
                    break;

                case 'fullScreen':
                    mainWin.setFullScreen(!mainWin.isFullScreen());
                    break;
                case 'minimize':
                    mainWin.minimize();
                    break;
                case 'maximize':
                    mainWin.maximize();
                    break;
                case 'unmaximize':
                    mainWin.unmaximize();
                    break;

                case 'help':
                    shell.openExternal('https://ztf.im');
                    break;
                case 'exit':
                    app.quit()
                    break;
                default:
                   logInfo('--', arg.action, arg.path)
                    if (arg.action == 'openInExplore')
                        this.openInExplore(arg.path)
                   else if (arg.action == 'openInTerminal')
                        this.openInTerminal(arg.path)
            }
        })
    }

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

         globalShortcut.register('CommandOrControl+D', () => {
             const mainWin = this._windows.get('main');
             mainWin.toggleDevTools()
         })
    }

    quit() {
        killZtfServer();
    }

    bindElectronEvents() {
        app.on('activate', () => {
            logInfo('>> event: activate');

            this.buildAppMenu();

            // 在 OS X 系统上，可能存在所有应用窗口关闭了，但是程序还没关闭，此时如果收到激活应用请求，需要重新打开应用窗口并创建应用菜单。
            this.openOrCreateWindow()
        });
        app.on('window-all-closed', () => {
            logInfo(`>> event: window-all-closed`)
            app.quit();
        });

        app.on('will-quit', (e) => {
            e.preventDefault()

            logInfo(`>> event: app will-quit`)

            logInfo(`>> start to kill child process`)
            this.quit();
            logInfo(`>> end to kill child process`)

            app.exit()
        });

        app.on('quit', () => {
            logInfo(`>> event: app quit`)
        });
    }

    showFileSelection(event) {
        dialog.showOpenDialog({
            properties: ['openFile']
        }).then(result => {
            this.reply(event, result)
        }).catch(err => {
            logErr(err)
        })
    }

    showFolderSelection(event) {
        dialog.showOpenDialog({
            properties: ['openDirectory']
        }).then(result => {
            this.reply(event, result)
        }).catch(err => {
            logErr(err)
        })
    }

    reply(event, result)  {
        if (result.filePaths && result.filePaths.length > 0) {
            event.reply(electronMsgReplay, result.filePaths[0]);
        }
    }

    openInExplore(path) {
        shell.showItemInFolder(path);
    }
    openInTerminal(path) {
        logInfo('openInTerminal')

        const stats = fs.statSync(path);
        if (stats.isFile()) {
            path = pth.resolve(path, '..')
        }

        if (IS_WINDOWS_OS) {
            cp.exec('start cmd.exe /K cd /D ' + path);
        } else if (IS_LINUX) {
            // support other terminal types
            cp.spawn ('gnome-terminal', [], { cwd: path });
        } else if (IS_MAC_OSX) {
            cp.exec('open -a Terminal ' + path);
        }
    }

    get windows() {

        return this._windows;
    }

    setAboutPanel() {
        if (!app.setAboutPanelOptions) {
            return;
        }

        let version = Config.pkg.buildTime ? 'build at ' + new Date(Config.pkg.buildTime).toLocaleString() : ''
        version +=  DEBUG ? '[debug]' : ''
        app.setAboutPanelOptions({
            applicationName: Lang.string(Config.pkg.name) || Config.pkg.displayName,
            applicationVersion: Config.pkg.version,
            copyright: `${Config.pkg.copyright} ${Config.pkg.company}`,
            credits: `Licence: ${Config.pkg.license}`,
            version: version
        });
    }

    buildAppMenu() {
        logInfo('>> ztf app: build application menu.');

        if (IS_MAC_OSX) {
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
                            label: 'Reload',
                            accelerator: 'Command+R',
                            click: () => {
                                this._windows.get('main').webContents.reload();
                            }
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
                            shell.openExternal('https://ztf.im');
                        }
                    }]
                }
            ];

            const menu = Menu.buildFromTemplate(template);
            Menu.setApplicationMenu(menu);
        } else {
            Menu.setApplicationMenu(null);
        }
    }
}
