<template>
  <div>
    <div class="head">
      <div class="title">配置列表</div>
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

  </div>
</template>

<script>

import {listConfig, removeConfig} from "../../../../api/manage";

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
  name: 'ConfigList',
  components: {
  },
  data() {
    return {
      models: [],
      columns,

      designVisible: false,
      designModel: {},
      time: 0,
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
      this.$router.push({path: '/data/buildin/config/edit/0'});
    },
    loadData() {
      listConfig().then(json => {
        console.log('listConfig', json)
        this.models = json.data
      })
    },
    edit(record) {
      console.log(record)
      this.$router.push({path: `/data/buildin/config/edit/${record.id}`});
    },
    remove(record) {
      console.log(record)
      removeConfig(record.id).then(json => {
        console.log('removeConfig', json)
        this.loadData()
      })
    },
  }
}
</script>

<style scoped>

</style>
