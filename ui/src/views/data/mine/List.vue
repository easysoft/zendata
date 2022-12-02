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
      <a-col :span="hasSelected ? 12 : 24">
        <div class="main-table">
          <div v-if="defs.length==0" class="no-data-tips">{{$t('tips.pls.refresh.data')}}</div>

          <template v-if="defs.length>0">
            <a-table :columns="columns" :data-source="defs" :pagination="false" rowKey="id" :custom-row="customRow">
              <a slot="recordTitle" slot-scope="text, record" @click="design(record)">{{record.title}}</a>

              <span slot="folderWithPath" slot-scope="text, record">
                <a-tooltip placement="top" overlayClassName="tooltip-light">
                  <template slot="title">
                    <span>{{record.path | replacePathSep}}</span>
                  </template>
                  <a>{{record.path | pathToRelated}}</a>
                </a-tooltip>
              </span>

              <span slot="action" slot-scope="record">
                <a @click="design(record)" :title="$t('action.design')"><Icon type="control" :style="{fontSize: '16px'}" /></a> &nbsp;
                <a @click="edit(record)" :title="$t('action.edit')"><Icon type="form" :style="{fontSize: '16px'}" /></a> &nbsp;
                <a @click="showDeleteConfirm(record)" :title="$t('action.delete')"><Icon type="delete" :style="{fontSize: '16px'}" /></a>
              </span>
            </a-table>

            <div class="pagination-wrapper">
              <a-pagination size="small" simple @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" />
            </div>
          </template>
        </div>
      </a-col>
      <a-col v-if="hasSelected" :span="12">
        <Preview :record="selectedRecord" />
      </a-col>
    </a-row>

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

    <a-modal
      :visible="editModalVisible"
      :title="editModalVisible ? editRecord ? `${$t('menu.data.edit')}: ${editRecord.title}` : $t('title.data.create') : ''"
      :footer="false"
      :centered="true"
      :width="700"
      @cancel="handleCancelEditModal"
    >
      <Edit
        :v-if="editModalVisible"
        :id="editModalVisible ? editRecord ? editRecord.id : 0 : null"
        :afterSave="handleEditSave"
      />
    </a-modal>
  </div>
</template>

<script>

import {Icon, Modal} from 'ant-design-vue'
import {listDef, removeDef} from "../../../api/manage";
import {DesignComponent} from '../../../components'
import {PageSize, ResTypeDef, replacePathSep, pathToRelated} from "../../../api/utils";
import debounce from "lodash.debounce"
import Preview from './Preview';
import Edit from './Edit';

export default {
  name: 'Mine',
  components: {
    DesignComponent,
    Icon,
    Preview,
    Edit,
  },
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
      defs: [],
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
    hasSelected: function() {
      if (!this.defs) return false

      return this.defs.some(x => x.id == this.selected);
    },
    selectedRecord: function() {
      if (!this.defs) return null

      return this.defs?.find(x => x.id == this.selected);
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
        console.log('listDefs', json)
        const that = this
        that.defs = json.data
        that.total = json.total
        that.selected = json.data.length ? json.data[0].id : null
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
