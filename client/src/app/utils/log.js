const logger = require('electron-log');
import {DEBUG} from './consts';
logger.transports.file.resolvePath = () =>
    require("path").join(require("os").homedir(), 'zd', 'log', 'electron.log');

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
