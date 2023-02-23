import path from "path";
import os from "os";
import {App, downloadUrl, ResDir, WorkDir} from "./consts";
import {app} from "electron";
import {logInfo} from "./log";

import {promisify} from 'node:util';
import stream from 'node:stream';
import fs from 'node:fs';
import got from 'got';
import crypto from "crypto";
const pipeline = promisify(stream.pipeline);

export function mkdir(dir) {
    const pth = path.join(WorkDir, dir);
    fs.mkdirSync(pth, {recursive:true})

    return pth
}

export function getResDir(dir) {
    const pth = path.resolve(process.resourcesPath, dir);
    return pth
}

export function getDownloadPath(version) {
    const pth = path.join(WorkDir, 'tmp', 'download', `${version}.zip`);
    return pth
}


export function getCurrVersion() {
    let currVersionStr = '0';

    const {versionPath} = getResPath()
    if (fs.existsSync(versionPath)) {
        const content = fs.readFileSync(versionPath)
        let json = JSON.parse(content);
        currVersionStr = json.version;
    } else {
        currVersionStr = app.getVersion();
    }

    const currVersion = parseFloat(currVersionStr);

    return {currVersion, currVersionStr};
}

export async function getRemoteVersion() {
    const versionUrl = getVersionUrl();

    const json = await got.get(versionUrl).json();
    const newVersionStr = json.version;
    const newVersion = parseFloat(newVersionStr);
    const forceUpdate = json.force;

    return {
        newVersion,
        newVersionStr,
        forceUpdate,
    }
}

export function changeVersion(newVersion) {
    const pth = path.resolve(ResDir, 'version.json');

    let json = {}
    if (fs.existsSync(pth)) {
        const content = fs.readFileSync(pth)
        json = JSON.parse(content);
    }

    json.version = newVersion;
    fs.writeFileSync(pth, JSON.stringify(json));
}

export function restart() {
    app.relaunch({
        args: process.argv.slice(1)
    });
    app.exit(0);
}

export function getResPath() {
    const versionPath = path.resolve(ResDir, 'version.json')
    const uiPath =  path.resolve(ResDir, 'ui');
    const serverPath = getBinPath('server')

    return {
        versionPath, uiPath, serverPath
    }
}

export function getBinPath(name) {
    const platform = os.platform(); // 'darwin', 'linux', 'win32'
    const execPath = `bin/${platform}/${name}${platform === 'win32' ? '.exe' : ''}`;
    const pth = path.join(ResDir, execPath);

    return pth
}

export function computerFileMd5(pth) {
    const buffer = fs.readFileSync(pth);
    const hash = crypto.createHash('md5');
    hash.update(buffer, 'utf8');
    const md5 = hash.digest('hex');
    return md5
}

export function getVersionUrl() {
    const url = new URL(`${App}/version.json`, downloadUrl) + '?ts=' + Date.now();
    logInfo(`versionUrl=${url}`)
    return url
}
export function getAppUrl(version) {
    const platform = os.platform(); // 'darwin', 'linux', 'win32'
    const url = new URL(`${App}/${version}/${platform}/${App}-upgrade.zip`, downloadUrl) + '?ts=' + Date.now();
    logInfo(`appUrl=${url}`)
    return url
}

export async function checkMd5(version, file) {
    const platform = os.platform(); // 'darwin', 'linux', 'win32'
    const url = new URL(`${App}/${version}/${platform}/${App}-upgrade.zip.md5`, downloadUrl);

    const md5Remote = await got.get(url).text();
    const md5File = computerFileMd5(file)
    const pass = md5Remote === md5File
    logInfo(`md5Remote=${md5Remote}, md5File=${md5File}, pass=${pass}`)

    return pass
}