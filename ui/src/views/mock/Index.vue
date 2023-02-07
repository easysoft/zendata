<template>
  <div class="mock-index-main">
    <div class="head">
      <div class="title">
        <Icon type="database" :style="{fontSize: '16px'}" />
        <span>{{$t('menu.data.mock')}}</span>
      </div>
      <div class="filter">
        <a-input-search v-model="keywords" @change="search" :allowClear="true" :placeholder="$t('tips.search')" style="width: 300px" />
      </div>
      <div class="buttons">
        <a-button type="primary" @click="create()"><Icon type="plus" :style="{fontSize: '16px'}" /> {{$t('action.create')}}</a-button>
      </div>
    </div>

    <a-row :gutter="10">
      <a-col :span="12">
        <List />
      </a-col>
      <a-col :span="12">
        <Preview />
      </a-col>
    </a-row>

  </div>
</template>

<script>

import {Icon, Modal} from 'ant-design-vue'
import {PageSize, ResTypeDef, replacePathSep, pathToRelated} from "../../api/utils";
import debounce from "lodash.debounce"
import Msg from '../../utils/msg.js'
import List from './List';
import Preview from './Preview';

export default {
  name: 'MockIndex',
  components: {
    Icon,
    List, Preview,
  },
  mixins: [],
  data() {
    return {
      keywords: '',
      page: 1,
      total: 0,
      pageSize: PageSize,
    };
  },

  methods: {
    search: debounce(function() {
      console.log('search', this.keywords)
      Msg.$emit('loadMock',{keywords: this.keywords})
    }, 500),

    create() {
      Msg.$emit('editMock',{})
    },
  }
}
</script>

<style lang="less" scoped>
  .mock-index-main {

  }
</style>
