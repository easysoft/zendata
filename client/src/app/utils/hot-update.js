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
import {execSync} from 'child_process';
import {IS_WINDOWS_OS} from "../utils/env";
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
            mainWin.webContents.send(electronMsgUpdate, {
                currVersionStr, newVersionStr, forceUpdate
            })
        }
    }
}

export const updateApp = (version, mainWin) => {
    downLoadApp(version, mainWin, doUpdate)
}

const doUpdate = async (downloadPath, version) => {
    await copyFiles(downloadPath);
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

    pipeline(downloadStream, fileWriterStream).then(async () => {
        logInfo(`success to downloaded to ${downloadPath}`)

        const md5Pass = await checkMd5(version, downloadPath)
        if (md5Pass) {
            cb(downloadPath, version)
        } else {
            logInfo('check md5 failed')
        }

    }).catch((err) => {
        logErr(`update failed: ${err}`)
    });
}

const copyFiles = async (downloadPath) => {
    const downloadDir = path.dirname(downloadPath)

    const extractedPath = path.resolve(downloadDir, 'extracted')
    logInfo(`downloadPath=${downloadPath}, extractedPath=${extractedPath}`)

    const unzip = new admZip(downloadPath, {});
    let pass = ''
    await unzip.extractAllTo(extractedPath, true, true, pass);
    logInfo(pass)

    const {uiPath, serverPath} = getResPath()
    logInfo(`uiPath=${uiPath}, serverPath=${serverPath}`)

    killZdServer();
    fs.rmSync(uiPath, {recursive: true})
    fs.rmSync(serverPath)

    const serverFile = `server${os.platform() === 'win32' ? '.exe' : ''}`
    const serverDir = path.dirname(serverPath)
    logInfo(`serverDir=${serverDir}`)
    const serverDist = path.join(serverDir, serverFile)
    logInfo(`serverFile=${serverFile}, serverDist=${serverDist}`)

    fse.copySync(path.resolve(downloadDir, 'extracted', 'ui'), path.resolve(path.dirname(uiPath), 'ui'), {recursive: true})
    fse.copySync(path.resolve(downloadDir, 'extracted', serverFile), serverDist)

    if (!IS_WINDOWS_OS) {
        const cmd = `chmod +x ${serverDist}`
        execSync(cmd, {windowsHide: true})
    }
}