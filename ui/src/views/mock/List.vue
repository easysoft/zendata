<template>
  <div class="mock-preview-list main-table">
    <div>
      <a-table :columns="columns" :data-source="models" :pagination="false" rowKey="id" :custom-row="customRow">
        <a slot="recordTitle" slot-scope="text, record" @click="view(record)">
          {{record.name}}
        </a>

        <a slot="createTime" slot-scope="text, record">
          {{record.createdAt | formatTime}}
        </a>

        <span slot="action" slot-scope="record">
                <a @click="edit(record)" :title="$t('action.edit')"><Icon type="form" :style="{fontSize: '16px'}" /></a> &nbsp;
                <a @click="showDeleteConfirm(record)" :title="$t('action.delete')">
                  <Icon type="delete" :style="{fontSize: '16px'}" />
                </a>
              </span>
      </a-table>

      <div class="pagination-wrapper">
        <a-pagination size="small" simple @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" />
      </div>
    </div>

    <div class="full-screen-modal">
      <mock-edit-comp
          ref="editComp"
          :type="type"
          :visible="editVisible"
          :model="editModel"
          :time="time"
          @ok="handleEditSave"
          @cancel="handleEditCancel" >
      </mock-edit-comp>
    </div>

  </div>
</template>

<script>

import {Icon, Modal} from 'ant-design-vue'
import {formatTime, PageSize, pathToRelated, replacePathSep, ResTypeDef} from "../../api/utils";
import debounce from "lodash.debounce"
import mockMixin from "@/store/mockMixin";
import Bus from '../../utils/bus.js'
import {listMock, removeMock} from "@/api/mock";
import MockEditComp from './components/Edit'

export default {
  name: 'MockList',
  components: {
    Icon,
    MockEditComp,
  },
  props: {
  },
  mixins: [mockMixin],
  filters: {
    formatTime: formatTime
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
        title: this.$i18n.t('msg.create.time'),
        dataIndex: 'createTime',
        scopedSlots: { customRender: 'createTime' },
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
      models: [],
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

      editVisible: false,
      editModel: null,
    };
  },
  computed: {
  },
  created () {
    this.loadData()
  },
  mounted () {
    Bus.$on('loadMock',(data) => {
      console.log('loadMock event', data)
      this.loadData()
    })

    Bus.$on('createMock',(data) => {
      console.log('createMock event', data)
      this.editModel = {}
      this.editVisible = true;
    })
  },
  methods: {
    loadData() {
      listMock(this.keywords, this.page).then(json => {
        this.models = json.data.list
        this.total = json.data.total
        this.selected = json.data.list.length ? json.data.list[0].id : null
      })
    },
    create() {
      this.editVisible = true;
    },
    edit(record) {
      this.editVisible = true;
      this.setMockItem(record)
    },
    handleEditSave() {
      this.editVisible = false;
      this.loadData();
    },
    handleEditCancel() {
      this.editVisible = false;
    },
    remove(record) {
      console.log(record)
      removeMock(record.id).then(json => {
        this.loadData()
      })
    },
    view(record) {
      this.setMockItem(record)
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
