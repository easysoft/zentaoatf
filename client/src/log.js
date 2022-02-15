const logger = require('electron-log');
logger.transports.file.resolvePath = () => require("path").join(require("os").homedir(), 'ztf', 'log', 'electron.log');

export function logInfo(str) {
    logger.info(str);
}
export function logErr(str) {
    logger.error(str);
}
