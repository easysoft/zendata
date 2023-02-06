<template>
  <div class="mock-preview-list main-table">
    <div v-if="mockItems.length==0" class="no-data-tips">{{$t('tips.pls.refresh.data')}}</div>

    <template v-if="mockItems.length>0">
      <a-table :columns="columns" :data-source="mockItems" :pagination="false" rowKey="id" :custom-row="customRow">
        <a slot="recordTitle" slot-scope="text, record" @click="view(record)">{{record.title}}</a>

        <span slot="folderWithPath" slot-scope="text, record">
                <a-tooltip placement="top" overlayClassName="tooltip-light">
                  <template slot="title">
                    <span>{{record.path | replacePathSep}}</span>
                  </template>
                  <a>{{record.path | pathToRelated}}</a>
                </a-tooltip>
              </span>

        <span slot="action" slot-scope="record">
                <a @click="edit(record)" :title="$t('action.edit')"><Icon type="form" :style="{fontSize: '16px'}" /></a> &nbsp;
                <a @click="showDeleteConfirm(record)" :title="$t('action.delete')"><Icon type="delete" :style="{fontSize: '16px'}" /></a>
              </span>
      </a-table>

      <div class="pagination-wrapper">
        <a-pagination size="small" simple @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" />
      </div>
    </template>
  </div>
</template>

<script>

import {Icon, Modal} from 'ant-design-vue'
import {listDef, removeDef} from "../../api/manage";
import {PageSize, ResTypeDef, replacePathSep, pathToRelated} from "../../api/utils";
import debounce from "lodash.debounce"
import mockMixin from "@/store/mockMixin";
import {listMock} from "@/api/mock";

export default {
  name: 'Mine',
  components: {
    Icon,
  },
  props: {
  },
  mixins: [mockMixin],
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
      mockItems: [],
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
    replacePathSep: function (path) {
      return replacePathSep(path)
    },
    pathToRelated: function (path) {
      return pathToRelated(path)
    }
  },
  methods: {
    loadData() {
      listMock(this.keywords, this.page).then(json => {
        this.mockItems = json.data.list
        this.total = json.data.total
        this.selected = json.data.list.length ? json.data.list[0].id : null
      })
    },
    create() {
      this.editModalVisible = true;
    },
    view(record) {
      this.setMockItem(record)
    },
    edit(record) {
      this.editModalVisible = true;
      this.setMockItem(record)
    },
    handleCancelEditModal() {
      this.editModalVisible = false;
    },
    handleEditSave() {
      this.editModalVisible = false;
      this.loadData();
    },

    remove(record) {
      console.log(record)
      removeDef(record.id).then(json => {
        console.log('removeDef', json)
        this.loadData()
      })
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
.mock-preview-list {
  .no-data-tips {
    padding: 15px;
    text-align: center;
  }
}
</style>
