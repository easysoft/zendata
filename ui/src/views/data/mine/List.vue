<template>
  <div>
    <div class="head">
      <div class="title">{{ $t('menu.data.list') }}</div>
      <div class="filter">
        <a-input-search v-model="keywords" @change="onSearch" :allowClear="true"
                        :placeholder="$t('tips.search')" style="width: 300px" />
      </div>
      <div class="buttons">
        <a-button type="primary" @click="create()">{{ $t('action.create') }}</a-button>
      </div>
    </div>

    <a-table :columns="columns" :data-source="defs" :pagination="false" rowKey="id">
      <span slot="folderWithPath" slot-scope="text, record">
        <a-tooltip placement="top" overlayClassName="tooltip-light">
          <template slot="title">
            <span>{{record.path}}</span>
          </template>
          <a>{{record.path | pathToRelated}}</a>
        </a-tooltip>
      </span>

      <span slot="action" slot-scope="record">
        <a @click="edit(record)">{{ $t('action.edit') }}</a> &nbsp;
        <a @click="design(record)">{{ $t('action.design') }}</a> &nbsp;

        <a-popconfirm
            :title="$t('tips.delete')"
            :okText="$t('msg.yes')"
            :cancelText="$t('msg.no')"
            @confirm="remove(record)"
          >
          <a href="#">{{ $t('action.delete') }}</a>
        </a-popconfirm> &nbsp;

        <a-popover :title="$t('msg.data')" @visibleChange="preview(record)" trigger="click"
                   placement="bottom" :autoAdjustOverflow="true">
          <template slot="content">
            <div v-html="previewData"></div>
          </template>
          <a>{{ $t('action.preview') }}</a>
        </a-popover>

      </span>
    </a-table>

    <div class="pagination-wrapper">
      <a-pagination @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" />
    </div>

    <div class="full-screen-modal">
      <design-component
          ref="designPage"
          :type="type"
          :visible="designVisible"
          :modelProp="designModel"
          :time="time"
          @ok="handleDesignOk"
          @cancel="handleDesignCancel" >
      </design-component>
    </div>

  </div>
</template>

<script>

import { listDef, removeDef, previewDefData } from "../../../api/manage";
import { DesignComponent } from '../../../components'
import {PageSize, ResTypeDef, pathToRelated} from "../../../api/utils";
import debounce from "lodash.debounce"

export default {
  name: 'Mine',
  components: {
    DesignComponent
  },
  data() {
    const columns = [
      {
        title: this.$i18n.t('form.name'),
        dataIndex: 'title',
      },
      {
        title: this.$i18n.t('form.file'),
        dataIndex: 'folder',
        scopedSlots: { customRender: 'folderWithPath' },
      },
      {
        title: this.$i18n.t('form.opt'),
        key: 'action',
        scopedSlots: { customRender: 'action' },
      },
    ];

    return {
      defs: [],
      previewData: '',
      columns,

      designVisible: false,
      designModel: {},
      type: ResTypeDef,
      time: 0,

      keywords: '',
      page: 1,
      total: 0,
      pageSize: PageSize,
    };
  },
  computed: {

  },
  created () {
    this.loadData()
  },
  mounted () {
  },
  filters: {
    pathToRelated: function (path) {
      return pathToRelated(path)
    }
  },
  methods: {
    loadData() {
      listDef(this.keywords, this.page).then(json => {
        console.log('listDefs', json)
        const that = this
        that.defs = json.data
        this.total = json.total
      })
    },
    create() {
      this.$router.push({path: '/data/mine/edit/0'});
    },
    edit(record) {
      console.log(record)
      this.$router.push({path: `/data/mine/edit/${record.id}`});
    },
    design(record) {
      this.time = Date.now() // trigger data refresh
      console.log(record)
      this.designVisible = true
      this.designModel = record
    },
    remove(record) {
      console.log(record)
      removeDef(record.id).then(json => {
        console.log('removeDef', json)
        this.loadData()
      })
    },
    preview(record) {
      console.log(record)
      previewDefData(record.id).then(json => {
        console.log('previewDefData', json)
        this.previewData = json.data
      })
    },

    handleDesignOk() {
      console.log('handleDesignOk')
      this.designVisible = false
    },
    handleDesignCancel() {
      console.log('handleDesignCancel')
      this.designVisible = false
      this.designModel = {}
    },

    onPageChange(page, pageSize) {
      console.log('onPageChange', page, pageSize)
      this.page= page
      this.loadData()
    },
    onSearch: debounce(function() {
      console.log('onSearch', this.keywords)
      this.loadData()
    }, 500),
  }
}
</script>

<style scoped>

</style>
