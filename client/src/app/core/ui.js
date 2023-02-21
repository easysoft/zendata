import {app} from 'electron';
import path from 'path';
import {spawn} from 'child_process';
import express from 'express';
import {portServer} from "../utils/consts";

import {DEBUG, portClient} from '../utils/consts';
import {logInfo, logErr} from '../utils/log';

let _uiService;

export function startUIService() {
    if (_uiService) {
        return Promise.resolve();
    }

    let {UI_SERVER_URL: uiServerUrl} = process.env;
    
    if (!uiServerUrl && !DEBUG) {
        return Promise.resolve("http://127.0.0.1:" + portServer + "/ui");
    }

    if (uiServerUrl) {
        if (/^https?:\/\//.test(uiServerUrl)) {
            return Promise.resolve(uiServerUrl);
        }
        return new Promise((resolve, reject) => {
            if (!path.isAbsolute(uiServerUrl)) {
                uiServerUrl = path.resolve(app.getAppPath(), uiServerUrl);
            }

            const port = process.env.UI_SERVER_PORT || portClient;
            logInfo(`>> starting ui serer at ${uiServerUrl} with port ${port}`);

            const uiServer = express();
            uiServer.use(express.static(uiServerUrl));
            const server = uiServer.listen(port, serverError => {
                if (serverError) {
                    console.error('>>> start ui server failed with error', serverError);
                    _uiService = null;
                    reject(serverError);
                } else {
                    logInfo(`>> ui server started successfully on http://localhost:${port}.`);
                    resolve(`http://localhost:${port}`);
                }
            });
            server.on('close', () => {
                _uiService = null;
            });
            _uiService = uiServer;
        })
    }

    return new Promise((resolve, reject) => {
        const cwd = path.resolve(app.getAppPath(), '../ui');
        logInfo(`>> starting ui development server with command "npm run serve" in "${cwd}"...`);

        let resolved = false;
        const cmd = spawn('npm', ['run', 'serve'], {
            cwd,
            shell: true,
        });
        cmd.on('close', (code) => {
            logInfo(`>> ui server closed with code ${code}`);
            _uiService = null;
        });
        cmd.stdout.on('data', data => {
            if (resolved) {
                return;
            }
            const dataString = String(data);
            const lines = dataString.split('\n');
            for (let i = 0; i < lines.length; i++) {
                const line = lines[i];
                if (DEBUG) {
                    logInfo('\t' + line);
                }
                if (line.includes('App running at:')) {
                    const nextLine = lines[i + 1] || lines[i + 2];
                    if (DEBUG) {
                        logInfo('\t' + nextLine);
                    }
                    if (!nextLine) {
                        logErr('\t' + `cannot grabing running address after line "${line}".`);
                        throw new Error(`cannot grabing running address after line "${line}".`);
                    }
                    const url = nextLine.split('Local:   ')[1];
                    if (url) {
                        resolved = true;
                        resolve(url);
                    }
                    if (!DEBUG) {
                        break;
                    }
                }
            }
        });
        cmd.on('error', spawnError => {
            logErr('>>> start ui server failed, error' + spawnError);
            reject(spawnError)
        });
        _uiService = cmd;
    });
}
