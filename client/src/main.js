import {app} from 'electron';
import {DEBUG} from './app/utils/consts';
import {ZtfApp} from "./app/app";
import {logInfo} from "./app/utils/log";

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) { // eslint-disable-line global-require
  app.quit();
}

logInfo(`DEBUG=${DEBUG}`)

const ztfApp = new ZtfApp();
app.on('ready', () => {
  ztfApp.ready()
});
