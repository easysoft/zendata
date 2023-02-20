/**
 * 格式化字符串
 * @param {string} str 要格式化的字符串
 * @param  {...any} args 格式化参数
 * @return  {string} 格式化后的字符串
 * @example <caption>通过参数序号格式化</caption>
 *     var hello = $.format('{0} {1}!', 'Hello', 'world');
 *     // hello 值为 'Hello world!'
 * @example <caption>通过对象名称格式化</caption>
 *     var say = $.format('Say {what} to {who}', {what: 'hello', who: 'you'});
 *     // say 值为 'Say hello to you'
 */
export const formatString = (str, ...args) => {
    let result = str;
    if (args.length > 0) {
        let reg;
        if (args.length === 1 && (typeof args[0] === 'object')) {
            args = args[0];
            Object.keys(args).forEach(key => {
                if (args[key] !== undefined) {
                    reg = new RegExp(`({${key}})`, 'g');
                    result = result.replace(reg, args[key]);
                }
            });
        } else {
            for (let i = 0; i < args.length; i++) {
                if (args[i] !== undefined) {
                    reg = new RegExp(`({[${i}]})`, 'g');
                    result = result.replace(reg, args[i]);
                }
            }
        }
    }
    return result;
};

/**
 * 字节单位表
 * @type {Object}
 */
export const BYTE_UNITS = {
    B: 1,
    KB: 1024,
    MB: 1024 * 1024,
    GB: 1024 * 1024 * 1024,
    TB: 1024 * 1024 * 1024 * 1024,
};

/**
 * 格式化字节值为包含单位的字符串
 * @param {number} size 字节大小
 * @param {number} [fixed=2] 保留的小数点尾数
 * @param {string} [unit=''] 单位，如果留空，则自动使用最合适的单位
 * @return {string} 格式化后的字符串
 */
export const formatBytes = (size, fixed = 2, unit = '') => {
    if (typeof size !== 'number') {
        size = Number.parseInt(size, 10);
    }
    if (Number.isNaN(size)) {
        return '?KB';
    }
    if (!unit) {
        if (size < BYTE_UNITS.KB) {
            unit = 'B';
        } else if (size < BYTE_UNITS.MB) {
            unit = 'KB';
        } else if (size < BYTE_UNITS.GB) {
            unit = 'MB';
        } else if (size < BYTE_UNITS.TB) {
            unit = 'GB';
        } else {
            unit = 'TB';
        }
    }

    return (size / BYTE_UNITS[unit]).toFixed(fixed) + unit;
};

/**
 * 检查字符串是否为未定义（`null` 或者 `undefined`）或者为空字符串
 * @param  {string} s 要检查的字符串
 * @return {boolean} 如果未定义或为空字符串则返回 `true`，否则返回 `false`
 */
export const isEmptyString = s => (s === undefined || s === null || s === '');

/**
 * 检查字符串是否不是空字符串
 * @param  {string} s 要检查的字符串
 * @return {boolean} 如果为非空字符串则返回 `true`，否则返回 `false`
 */
export const isNotEmptyString = s => (s !== undefined && s !== null && s !== '');

/**
 * 检查字符串是否不是空字符串，如果为空则返回第二个参数给定的字符串，否则返回字符串自身
 * @param  {string} str 要检查的字符串
 * @param  {string} thenStr 如果为空字符串时要返回的字符串
 * @return {boolean} 如果未定义或为空字符串则返回 [thenStr]，否则返回 [str]
 */
export const ifEmptyStringThen = (str, thenStr) => (isEmptyString(str) ? thenStr : str);

/**
 * 确保字符串长度不超过指定值，如果超出则去掉截取的部分
 * @param {string} str 要操作的字符串
 * @param {number} length 要限制的最大长度
 * @param {string} suffix 如果超出显示要添加的后缀
 * @returns {string} 返回新的字符串
 */
export const limitStringLength = (str, length, suffix) => {
    if (str.length > length) {
        str = str.substring(0, length);
        if (suffix) {
            str = `${str}${suffix}`;
        }
    }
    return str;
};

/**
 * 用于匹配 @ 用户的正则表达式
 * @type {string}
 */
export const REGEXP_AT_USER = '@(#?[_.\\w\\d\\u4e00-\\u9fa5]{1,20})';

/**
 * 还原包含 @ 成员的文本消息
 * @param {string} message @ 成员消息
 * @returns {string} 原是文本
 */
export const restoreMessageContainAt = (message) => {
    if (typeof message !== 'string' || !message.replace) {
        return message;
    }
    return message.replace(new RegExp(`\\[(?<atuser>${REGEXP_AT_USER})\\]\\(\\@\\#\\d+\\)`, 'g'), '$<atuser>');
};

export default {
    format: formatString,
    isEmpty: isEmptyString,
    isNotEmpty: isNotEmptyString,
    formatBytes,
    ifEmptyThen: ifEmptyStringThen,
    limitLength: limitStringLength
};
