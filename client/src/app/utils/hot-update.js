import {app} from 'electron';
import {
    changeVersion, checkMd5, getAppUrl,
    getCurrVersion,
    getDownloadPath,
    getRemoteVersion,
    getResPath, mkdir,
    restart
} from "./comm";
import {electronMsgDownloading, electronMsgUpdate, WorkDir} from "./consts";
import path from "path";
import fse from 'fs-extra'
import {logErr, logInfo} from "./log";

export async function checkUpdate(mainWin) {
    logInfo('checkUpdate ...')

    const {currVersion, currVersionStr} = getCurrVersion()

    const {newVersion, newVersionStr, forceUpdate} = await getRemoteVersion()
    logInfo(`currVersion=${currVersion}, newVersion=${newVersion}, forceUpdate=${forceUpdate}`)
    logInfo(currVersion < newVersion)
    if (currVersion < newVersion) {
        if (forceUpdate) {
            // logInfo('forceUpdate')
        } else {
            logInfo(11)
            mainWin.webContents.send(electronMsgUpdate, {
                currVersionStr, newVersionStr, forceUpdate
            })
        }
    }
}

export const updateApp = (version, mainWin) => {
    downLoadApp(version, mainWin, doUpdate)
}

const doUpdate = (downloadPath, version)=>{
    copyFiles(downloadPath);
    changeVersion(version);
    restart();
}

const admZip = require('adm-zip');

import {promisify} from 'node:util';
import stream from 'node:stream';
import fs from 'node:fs';
import got from 'got';
import {cpSync} from "fs";
import os from "os";
import {killZdServer} from "../core/zd";
const pipeline = promisify(stream.pipeline);

mkdir(path.join('tmp', 'download'))

const downLoadApp = (version, mainWin, cb) => {
    const downloadUrl = getAppUrl(version)
    const downloadPath = getDownloadPath(version)

    const downloadStream = got.stream(downloadUrl);
    const fileWriterStream = fs.createWriteStream(downloadPath);

    logInfo(`start download ${downloadUrl} ...`)

    downloadStream.on("downloadProgress", ({ transferred, total, percent }) => {
        mainWin.webContents.send(electronMsgDownloading, {percent})
    });

    pipeline(downloadStream, fileWriterStream).then(() => {
        logInfo(`success to downloaded to ${downloadPath}`)

        const md5Pass = checkMd5(version, downloadPath)
        if (md5Pass) {
            cb(downloadPath, version)
        } else {
            logInfo('check md5 failed')
        }

    }).catch((err) => {
        logErr(`update failed: ${err}`)
    });
}

const copyFiles = (downloadPath) => {
    const downloadDir = path.dirname(downloadPath)

    const extractedPath = path.resolve(downloadDir, 'extracted')
    logInfo(`downloadPath=${downloadPath}, extractedPath=${extractedPath}`)

    const unzip = new admZip(downloadPath);
    unzip.extractAllTo(extractedPath, true);

    const {uiPath, serverPath} = getResPath()
    logInfo(`uiPath=${uiPath}, serverPath=${serverPath}`)

    killZdServer();
    fs.rmSync(serverPath)
    fs.rmSync(uiPath, { recursive:true })

    const serverFile = `server${os.platform() === 'win32' ? '.exe' : ''}`
    fse.copySync(path.resolve(downloadDir, 'extracted', 'ui'), path.resolve(path.dirname(uiPath), 'ui'), { recursive: true })
    fse.copySync(path.resolve(downloadDir, 'extracted', serverFile), path.resolve(path.dirname(serverPath), serverFile))
}