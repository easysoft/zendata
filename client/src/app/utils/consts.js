import path from "path";
import os from "os";

export const DEBUG = true; // process.env.NODE_ENV === 'development';
export const WORK_DIR = process.cwd()

export const portClient = 55233
export const portServer = 55234
export const uuid = 'ZENDATA@1CF17A46-B136-4AEB-96B4-F21C8200EF5A~'

export const electronMsg = 'electronMsg'
export const electronMsgReplay = 'electronMsgReplay'
export const electronMsgUpdate = 'electronMsgUpdate'
export const electronMsgDownloading = 'electronMsgDownloading'

export const minimumSizeWidth = 1024
export const minimumSizeHeight = 640

export const App = 'zd';
export const WorkDir = path.join(os.homedir(), App);
export const ResDir = process.resourcesPath;
export const downloadUrl = 'https://dl.cnezsoft.com/';