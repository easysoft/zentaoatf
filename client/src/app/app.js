import {app, BrowserWindow} from 'electron';

import Lang from './utils/lang';
import {IS_MAC_OSX, setEntryPath} from './utils/env';
import {logInfo} from './utils/log';
import {getUIServerUrl, killZtfServer} from "./service";

export default class ZtfApp {
    constructor(entryPath) {
        setEntryPath(entryPath);

        this.lang = Lang;
        this.bindElectronEvents();

        logInfo(`>> ZtfApp: created, entry path is "${entryPath}".`);
    }

    ready() {
        logInfo('>> ZtfApp: ready.');
        this.openOrCreateWindow();
        // this.buildAppMenu();
    }

    openOrCreateWindow() {
        const {lastFocusedAppWin} = this;
        if (lastFocusedAppWin) {
            lastFocusedAppWin.showAndFocus();
        } else {
            this.createWindow();
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

        getUIServerUrl().then((url) => {
            mainWindow.loadURL(url);
            mainWindow.webContents.openDevTools({mode: 'bottom'});
        })
    };

    bindElectronEvents() {
        app.on('window-all-closed', () => {
            logInfo(`>> Event window-all-closed`)
            app.quit();
        });

        app.on('quit', () => {
            logInfo(`>> Event quit`)
            killZtfServer();
        });

        app.on('activate', () => {
            logInfo('>> ElectronApp: activate');

            // 在 OS X 系统上，可能存在所有应用窗口关闭了，但是程序还没关闭，此时如果收到激活应用请求需要
            // 重新打开应用窗口并创建应用菜单
            this.createWindow();
            // this.buildAppMenu();
        });
    }

    // buildAppMenu() {
    //     if (!IS_MAC_OSX) {
    //         return;
    //     }
    //
    //     const template = [{
    //         label: Lang.string('app.title', Config.pkg.displayName),
    //         submenu: [{
    //             label: Lang.string('menu.about'),
    //             selector: 'orderFrontStandardAboutPanel:'
    //         }, {
    //             type: 'separator'
    //         }, {
    //             label: 'Services',
    //             submenu: []
    //         }, {
    //             type: 'separator'
    //         }, {
    //             label: Lang.string('menu.hideCurrentWindow'),
    //             accelerator: 'Command+H',
    //             selector: 'hide:'
    //         }, {
    //             label: Lang.string('menu.hideOtherWindows'),
    //             accelerator: 'Command+Shift+H',
    //             selector: 'hideOtherApplications:'
    //         }, {
    //             label: Lang.string('menu.showAllWindows'),
    //             selector: 'unhideAllApplications:'
    //         }, {
    //             type: 'separator'
    //         }, {
    //             label: Lang.string('menu.quit'),
    //             accelerator: 'Command+Q',
    //             click: () => {
    //                 this.quit();
    //             }
    //         }]
    //     },
    //         {
    //             label: Lang.string('menu.edit'),
    //             submenu: [{
    //                 label: Lang.string('menu.undo'),
    //                 accelerator: 'Command+Z',
    //                 selector: 'undo:'
    //             }, {
    //                 label: Lang.string('menu.redo'),
    //                 accelerator: 'Shift+Command+Z',
    //                 selector: 'redo:'
    //             }, {
    //                 type: 'separator'
    //             }, {
    //                 label: Lang.string('menu.cut'),
    //                 accelerator: 'Command+X',
    //                 selector: 'cut:'
    //             }, {
    //                 label: Lang.string('menu.copy'),
    //                 accelerator: 'Command+C',
    //                 selector: 'copy:'
    //             }, {
    //                 label: Lang.string('menu.paste'),
    //                 accelerator: 'Command+V',
    //                 selector: 'paste:'
    //             }, {
    //                 label: Lang.string('menu.selectAll'),
    //                 accelerator: 'Command+A',
    //                 selector: 'selectAll:'
    //             }]
    //         },
    //         {
    //             label: Lang.string('menu.view'),
    //             submenu: (DEBUG) ? [{
    //                 label: Lang.string('menu.reload'),
    //                 accelerator: 'Command+R',
    //                 click: () => {
    //                     this.getLastFocusedWindow({includeChildWindow: true}).browserWindow.webContents.reload();
    //                 }
    //             }, {
    //                 label: Lang.string('menu.toggleFullscreen'),
    //                 accelerator: 'Ctrl+Command+F',
    //                 click: () => {
    //                     const lastWindow = this.getLastFocusedWindow();
    //                     lastWindow.browserWindow.setFullScreen(!lastWindow.browserWindow.isFullScreen());
    //                 }
    //             }, {
    //                 label: Lang.string('menu.toggleDeveloperTool'),
    //                 accelerator: 'Alt+Command+I',
    //                 click: () => {
    //                     this.lastFocusedAppWin.browserWindow.toggleDevTools();
    //                 }
    //             }] : [{
    //                 label: Lang.string('menu.toggleFullscreen'),
    //                 accelerator: 'Ctrl+Command+F',
    //                 click: () => {
    //                     const lastWindow = this.getLastFocusedWindow();
    //                     lastWindow.browserWindow.setFullScreen(!lastWindow.browserWindow.isFullScreen());
    //                 }
    //             }]
    //         },
    //         {
    //             label: Lang.string('menu.window'),
    //             submenu: [{
    //                 label: Lang.string('menu.createNewWindow'),
    //                 accelerator: 'Command+N',
    //                 click: () => {
    //                     this.createMainWindow();
    //                 }
    //             }, {
    //                 label: Lang.string('menu.minimize'),
    //                 accelerator: 'Command+M',
    //                 selector: 'performMiniaturize:'
    //             }, {
    //                 label: Lang.string('menu.close'),
    //                 accelerator: 'Command+W',
    //                 selector: 'performClose:'
    //             }, {
    //                 type: 'separator'
    //             }, {
    //                 label: Lang.string('menu.bringAllToFront'),
    //                 selector: 'arrangeInFront:'
    //             }]
    //         },
    //         {
    //             label: Lang.string('menu.help'),
    //             submenu: [{
    //                 label: Lang.string('menu.website'),
    //                 click: () => {
    //                     shell.openExternal(Lang.string('app.homepage', Config.pkg.homepage));
    //                 }
    //             }, {
    //                 label: Lang.string('menu.community'),
    //                 click() {
    //                     shell.openExternal('https://www.xuanim.com/forum/');
    //                 }
    //             }]
    //         }];
    //
    //     const menu = Menu.buildFromTemplate(template);
    //     Menu.setApplicationMenu(menu);
    //
    //     if (DEBUG) {
    //         console.log('>> XuanxuanApp: build application menu.');
    //     }
    // }
}