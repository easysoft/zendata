const logger = require('electron-log');
const path = require("path")
import {DEBUG} from './consts';

export const logDir = path.join(require("os").homedir(), 'zd', 'log')
logger.transports.file.resolvePath = () => path.join(logDir, 'electron.log');

export function logDebug(str) {
    if (DEBUG) {
        logger.debug(str);
    }
}
export function logInfo(str) {
    logger.info(str);
}
export function logErr(str) {
    logger.error(str);
}
