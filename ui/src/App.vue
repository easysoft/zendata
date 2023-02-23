<template>
  <a-config-provider :locale="locale">
    <div id="app">
      <router-view />
    </div>
  </a-config-provider>
</template>

<script>

import { domTitle, setDocumentTitle } from './utils/dom'
import { i18nRender } from '@/locales'

export default {
  name: 'App',
  components: {
  },

  data() {
    return {
    };
  },
  computed: {
    locale () {
      // 切换语言时，更新标题
      const { title } = this.$route.meta
      title && (setDocumentTitle(`${i18nRender(title)} - ${domTitle}`))
      return this.$i18n.getLocaleMessage(this.$store.getters.lang).antLocale
    }
  }
}
</script>

<style>
#app {
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}
</style>
