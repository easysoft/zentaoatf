import {app, BrowserWindow} from 'electron';
import {getUIServerUrl, startZtfServer, killZtfServer} from './service';
import {logInfo, logErr} from './log';

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) { // eslint-disable-line global-require
  app.quit();
}

const createWindow = (url) => {
  // Create the browser window.
  // const mainWindow = new BrowserWindow({
  //   width: 1200,
  //   height: 800,
  // });

  require('@electron/remote/main').initialize()

  const mainWindow = new BrowserWindow({show: false, webPreferences: {nodeIntegration: true, contextIsolation: false}})

  require("@electron/remote/main").enable(mainWindow.webContents)

  mainWindow.maximize()
  mainWindow.show()

  // and load the index.html of the app.
  mainWindow.loadURL(url);

  // Open the DevTools.
  mainWindow.webContents.openDevTools({mode: 'bottom'});
};

let _starting = false;

async function startApp() {
  if (_starting) {
    return;
  }
  _starting = true;

  // try {
  //   const ztfServerUrl = await startZtfServer();
  //   logInfo(`>> ZTF Server started successfully: ${ztfServerUrl}`);
  // } catch (error) {
  //   logErr('>> Start ztf server failed: ' + error);
  //   process.exit(1);
  //   return;
  // }

  const url = await getUIServerUrl();

  logInfo(`>> UI server url is ${url}`);

  createWindow(url);

  _starting = false;
}

// This method will be called when Electron has finished
// initialization and is ready to create browser windows.
// Some APIs can only be used after this event occurs.
app.on('ready', startApp);

// Quit when all windows are closed, except on macOS. There, it's common
// for applications and their menu bar to stay active until the user quits
// explicitly with Cmd + Q.
app.on('window-all-closed', () => {
  logInfo(`>> Event window-all-closed`)

  // if (process.platform !== 'darwin') {
  app.quit();
  // }
});

app.on('before-quit',function(){
  logInfo(`>> Event before-quit`)
})

app.on('will-quit',function(){
  logInfo(`>> Event will-quit`)
})

app.on('quit',function(){
  logInfo(`>> Event quit`)
  killZtfServer();
})

app.on('activate', () => {
  logInfo(`>> Event activate`)

  // On OS X it's common to re-create a window in the app when the
  // dock icon is clicked and there are no other windows open.
  if (BrowserWindow.getAllWindows().length === 0) {
    startApp();
  }
});
