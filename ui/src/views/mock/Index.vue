<template>
  <div class="container">
    <div class="head">
      <div class="title"><Icon type="database" :style="{fontSize: '16px'}" /><span>{{$t('menu.data.list')}}</span></div>
      <div class="filter">
        <a-input-search v-model="keywords" @change="onSearch" :allowClear="true" :placeholder="$t('tips.search')" style="width: 300px" />
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
import {listDef, removeDef} from "../../api/manage";
import {PageSize, ResTypeDef, replacePathSep, pathToRelated} from "../../api/utils";
import debounce from "lodash.debounce"
import List from './List';
import Preview from './Preview';
import mockMixin from "@/store/mockMixin";

export default {
  name: 'Mine',
  components: {
    Icon,
    List, Preview,
  },
  mixins: [],
  data() {
    const columns = [
      {
        title: this.$i18n.t('form.name'),
        dataIndex: 'title',
        'class': 'title',
        scopedSlots: { customRender: 'recordTitle' },
      },
      {
        title: this.$i18n.t('form.file'),
        dataIndex: 'folder',
        scopedSlots: { customRender: 'folderWithPath' },
        width: '300px'
      },
      {
        title: this.$i18n.t('form.opt'),
        key: 'action',
        scopedSlots: { customRender: 'action' },
        width: '80px'
      },
    ];

    return {
      columns,
      selected: null,

      designVisible: false,
      designModel: {},
      type: ResTypeDef,
      time: 0,

      keywords: '',
      page: 1,
      total: 0,
      pageSize: PageSize,

      editModalVisible: false,
      editRecord: null,
    };
  },
  computed: {
    // hasSelected: function() {
    //   if (!this.records) return false
    //
    //   return this.records.some(x => x.id == this.selected);
    // },
    selectedRecord: function() {
      return this.mockItem
    }
  },
  created () {
    this.loadData()
  },
  mounted () {
  },
  filters: {
    replacePathSep: function (path) {
      return replacePathSep(path)
    },
    pathToRelated: function (path) {
      return pathToRelated(path)
    }
  },
  methods: {
    loadData() {
      listDef(this.keywords, this.page).then(json => {
        this.records = json.data.list
        this.total = json.data.total
        this.selected = json.data.list.length ? json.data.list[0].id : null
      })
    },
    create() {
      this.editRecord = {};
      this.editModalVisible = true;
    },
    edit(record) {
      this.editRecord = record;
      this.editModalVisible = true;
    },
    handleCancelEditModal() {
      this.editModalVisible = false;
    },
    handleEditSave() {
      this.editModalVisible = false;
      this.loadData();
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
    handleClickRow: function(event) {
      const id = event.target.closest('tr').getAttribute('data-row-key');
      this.selected = id;
    },
    customRow: function(record) {
      const {selected} = this;
      return {
        attrs: {
          'class': record.id == selected ? 'selected' : ''
        },
        on: {
          click: this.handleClickRow
        }
      }
    },
    showDeleteConfirm: function(record) {
      Modal.confirm({
        title: this.$t('tips.delete'),
        content: (h) => <strong>{record.title}</strong>,
        okText: this.$t('msg.yes'),
        cancelText: this.$t('msg.no'),
        cancelType: 'danger',
        onOk: () => {
          this.remove(record)
        },
      });
    }
  }
}
</script>

<style lang="less" scoped>
.no-data-tips {
  padding: 15px;
  text-align: center;
}
</style>
