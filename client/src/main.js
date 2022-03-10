import {app} from 'electron';
import {DEBUG} from './app/utils/consts';
import ZtfApp from "./app/app";

// Handle creating/removing shortcuts on Windows when installing/uninstalling.
if (require('electron-squirrel-startup')) { // eslint-disable-line global-require
  app.quit();
}

console.log(`DEBUG=${DEBUG}`)

const ztfApp = new ZtfApp(__dirname);
app.on('ready', () => {
  ztfApp.ready()
});
