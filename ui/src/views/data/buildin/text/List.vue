<template>
  <div>
    <div class="head">
      <div class="title">文本列表</div>
      <div class="filter">
        <a-input-search v-model="keywords" @change="onSearch" :allowClear="true" placeholder="输入关键字检索" style="width: 300px" />
      </div>
      <div class="buttons">
        <a-button type="primary" @click="create()">新建</a-button>
      </div>
    </div>

    <a-table :columns="columns" :data-source="models" rowKey="id">
      <span slot="folderWithPath" slot-scope="text, record">
        <a-tooltip placement="top" overlayClassName="tooltip-light">
          <template slot="title">
            <span>{{record.path}}</span>
          </template>
          {{record.folder}}
        </a-tooltip>
      </span>

      <span slot="action" slot-scope="record">
        <a @click="edit(record)">编辑</a> |

        <a-popconfirm
            title="确认删除？"
            ok-text="是"
            cancel-text="否"
            @confirm="remove(record)"
        >
          <a href="#">删除</a>
        </a-popconfirm>
      </span>
    </a-table>

    <div class="pagination-wrapper">
      <a-pagination @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" />
    </div>

  </div>
</template>

<script>

import {listText, removeText} from "../../../../api/manage";
import {PageSize} from "../../../../api/utils";
import debounce from "lodash.debounce"

const columns = [
  {
    title: '名称',
    dataIndex: 'title',
  },
  {
    title: '引用',
    dataIndex: 'name',
  },
  {
    title: '目录',
    dataIndex: 'folder',
    scopedSlots: { customRender: 'folderWithPath' },
  },
  {
    title: '操作',
    key: 'action',
    scopedSlots: { customRender: 'action' },
  },
];

export default {
  name: 'TextList',
  components: {
  },
  data() {
    return {
      models: [],
      columns,

      designVisible: false,
      designModel: {},
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
  methods: {
    create() {
      this.$router.push({path: '/data/buildin/text/edit/0'});
    },
    loadData() {
      listText(this.keywords, this.page).then(json => {
        console.log('listText', json)
        this.models = json.data
      })
    },
    edit(record) {
      console.log(record)
      this.$router.push({path: `/data/buildin/text/edit/${record.id}`});
    },
    remove(record) {
      console.log(record)
      removeText(record.id).then(json => {
        console.log('removeText', json)
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
    },

    onPageChange() {
      console.log('onPageChange')
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
