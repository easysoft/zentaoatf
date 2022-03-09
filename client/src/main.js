import {app} from 'electron';
import ZtfApp from "./app/app";

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) { // eslint-disable-line global-require
  app.quit();
}

const ztfApp = new ZtfApp(__dirname);
app.on('ready', () => {
  ztfApp.ready()
});
