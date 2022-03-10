/**
 * 是否是 Mac OS 系统
 * @type {boolean}
 * @private
 */
export const IS_MAC_OSX = process.platform === 'darwin';

/**
 * 是否是 Windows 系统
 * @type {boolean}
 * @private
 */
export const IS_WINDOWS_OS = process.platform === 'win32';
/**
 * 是否是 Linux 系统
 * @type {boolean}
 * @private
 */
export const IS_LINUX = process.platform === 'linux';

/**
 * 获取Electron 应用入口文件所在目录
 * @returns {string} 目录路径
 */
export function getEntryPath() {
    return global.entryPath;
}

/**
 * 获取Electron 应用入口文件所在目录
 * @param {string} entryPath 目录路径
 * @return {void}
 */
export function setEntryPath(entryPath) {
    global.entryPath = entryPath;
}
