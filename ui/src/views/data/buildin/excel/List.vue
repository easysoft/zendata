<template>
  <div>
    <div class="head">
      <div class="title">表格列表</div>
      <div class="filter">
        <a-input-search v-model="keywords" @change="onSearch" :allowClear="true" placeholder="输入关键字检索" style="width: 300px" />
      </div>
      <div class="buttons">
        <a-button type="primary" @click="create()">新建</a-button>
      </div>
    </div>

    <a-table :columns="columns" :data-source="models" :pagination="false" rowKey="id">
      <span slot="folderWithPath" slot-scope="text, record">
        <a-tooltip placement="top" overlayClassName="tooltip-light">
          <template slot="title">
            <span>{{record.path}}</span>
          </template>
          <a>{{record.path | pathToRelated}}</a>
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
        </a-popconfirm> |

        <a-tooltip placement="top" overlayClassName="tooltip-light">
          <template slot="title">
            <div class="content-width" style="min-width: 280px;">
              <div class="title">引用Excel表格</div>
              <div class="content">
                <div>from: {{ record.referName }}</div>
                <div>select: field_name</div>
                <div>where: field_name like '%keywords%'</div>
              </div>
            </div>
          </template>
          <a href="#">引用</a>
        </a-tooltip>

      </span>

    </a-table>

    <div class="pagination-wrapper">
      <a-pagination @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" />
    </div>

  </div>
</template>

<script>

import {listExcel, removeExcel} from "../../../../api/manage";
import {PageSize, pathToRelated} from "../../../../api/utils";
import debounce from "lodash.debounce"

const columns = [
  {
    title: '名称',
    dataIndex: 'title',
  },
  {
    title: '文件',
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
  name: 'ExcelList',
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
  filters: {
    pathToRelated: function (path) {
      return pathToRelated(path)
    }
  },
  methods: {
    create() {
      this.$router.push({path: '/data/buildin/excel/edit/0'});
    },
    loadData() {
      listExcel(this.keywords, this.page).then(json => {
        console.log('listExcel', json)
        this.models = json.data
        this.total = json.total
      })
    },
    edit(record) {
      console.log(record)
      this.$router.push({path: `/data/buildin/excel/edit/${record.id}`});
    },
    remove(record) {
      console.log(record)
      removeExcel(record.id).then(json => {
        console.log('removeExcel', json)
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
