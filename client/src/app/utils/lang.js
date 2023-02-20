import TextMap from './text-map';
import {formatString} from './string';

export class LangHelper extends TextMap {
    constructor(name, langData) {
        super(langData);
        this._name = name;
    }

    change(name, langData) {
        this._data = langData;
        this._name = name;
    }

    get name() {
        return this._name;
    }

    error(err) {
        if (!err) {
            if (DEBUG) {
                console.collapse('LANG.error', 'redBg', '<Unknown Error>', 'redPale');
                console.error(err);
                console.groupEnd();
            }
            return '<Unknown Error>';
        }
        if (typeof err === 'string') {
            return this.string(err.startsWith('error.') ? err : `error.${err}`, err);
        }
        if (Array.isArray(err)) {
            return err.map(this.error).join(';');
        }
        let message = '';
        if (err.code) {
            message += this.string(`error.${err.code}`, `${err.message || ''}[${err.code}]`);
        } else if (err.message) {
            message = this.string(`error.${err.message}`, err.message);
        }
        if (message) {
            let formatParams = err.formats || err.extras;
            if (formatParams) {
                if (typeof formatParams === 'object' && !Array.isArray(formatParams)) {
                    message = formatString(message, formatParams);
                } else {
                    if (!Array.isArray(formatParams)) {
                        formatParams = [formatParams];
                    }
                    message = formatString(message, ...formatParams);
                }
            }
        }
        return message;
    }
}
