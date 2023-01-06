<template>
  <div class="container">
    <div class="head">
      <a-menu
        :default-selected-keys="['config']"
        :selectedKeys="[selectedKey]"
        mode="horizontal"
        @click="handleMenuClick"
        class="navbar-secondary"
      >
        <a-menu-item key="config">
          {{ $t('msg.config') }}
        </a-menu-item>
        <a-menu-item key="ranges">
          {{ $t('msg.ranges') }}
        </a-menu-item>
        <a-menu-item key="instances">
          {{ $t('msg.instances') }}
        </a-menu-item>
        <a-menu-item key="text">
          {{ $t('msg.text') }}
        </a-menu-item>
        <a-menu-item key="excel">
          {{ $t('msg.excel') }}
        </a-menu-item>
      </a-menu>
      <div class="filter">
        <a-input-search v-model="keywords" @change="onSearch" :allowClear="true" :placeholder="$t('tips.search')" style="width: 300px" />
      </div>

      <div class="buttons">
        <a-button v-if="selectedKey !== 'excel'" type="primary" @click="handleCreateClick()">
          <a-icon type="plus" :style="{fontSize: '16px'}" />
          {{$t('action.create')}}
        </a-button>
      </div>

    </div>
    <router-view />
  </div>
</template>

<script>
export default {
  name: 'BuildinLayout',
  data () {
    return {
      selected: this.$route.path.split('/')[3] || 'config',
      keywords: '',
      createShow: false
    }
  },
  computed: {
    selectedKey: function() {
      return this.$route.path.split('/')[3];
    }
  },
  methods: {
    handleMenuClick: function(e) {
      this.selected = e.key;
      this.keywords = '';
      this.createShow = false;
      this.updateRoutePath();
    },
    onSearch: function() {
      this.updateRoutePath();
    },
    updateRoutePath: function() {
      const {selected, keywords, createShow} = this;
      const path = `/data/buildin/${selected}/list${createShow ? '/0' : ''}`;
      const {query = {}} = this.$router;
      const oldKeywords = typeof query.search === 'string' ? query.search : '';
      if (this.$route.path !== path || oldKeywords !== keywords) {
        if (keywords.length) {
          query.search = keywords;
        }
        this.$router.push({
          path,
          query
        }).then(() => {
          console.log(this.$route);
        });
      }
    },
    handleCreateClick: function() {
      this.createShow = true;
      this.updateRoutePath();
    }
  },
}
</script>
