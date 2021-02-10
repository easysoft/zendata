import { i18nRender } from '@/locales'
import i18nMixin from '../../store/i18nMixin'
import { Dropdown, Icon, Menu } from 'ant-design-vue'
import './index.less'

const locales = ['zh-CN', 'en-US']
const languageLabels = {
  'zh-CN': '简体中文',
  'en-US': 'English'
}
const languageIcons = {
  'zh-CN': '🇨🇳',
  'en-US': '🇺🇸'
}

const SelectLang = {
  props: {
    prefixCls: {
      type: String,
    }
  },
  name: 'SelectLang',
  mixins: [i18nMixin],
  render () {
    const { prefixCls } = this
    const changeLang = ({ key }) => {
      this.setLang(key)
    }
    const langMenu = (
      <Menu class={['menu', 'ant-pro-header-menu']} selectedKeys={[this.currentLang]} onClick={changeLang}>
        {locales.map(locale => (
          <Menu.Item key={locale}>
            <span role="img" aria-label={languageLabels[locale]}>
              {languageIcons[locale]}
            </span>{' '}
            {languageLabels[locale]}
          </Menu.Item>
        ))}
      </Menu>
    )
    return (
      <Dropdown overlay={langMenu} placement="bottomRight">
        <span class={prefixCls}>
          <Icon type="global" title={i18nRender('navBar.lang')} style={{fontSize: '18px'}} />
        </span>
      </Dropdown>
    )
  }
}

export default SelectLang
