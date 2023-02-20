import pkg from '../package.json';

/**
 * 运行时配置对象
 * @type {{pkg: Object, media: Object, system: Object}}
 */
const config = {pkg};

/**
 * 更新运行时配置对象
 * @param {Object} newConfig 新的配置对象
 * @return {void}
 */
export function updateConfig(newConfig) {
    Object.assign(config, newConfig);
}

export default config;
