<template>
  <div>
    <div class="head">
      <div class="title">测试数据列表</div>
      <div class="buttons">
        <a-button type="primary" @click="create()">新建</a-button>
      </div>
    </div>

    <a-table :columns="columns" :data-source="defs" rowKey="id">
      <a slot="name" slot-scope="text">{{ text }}</a>

      <span slot="customTitle">名称</span>

      <span slot="action" slot-scope="record">
        <a @click="edit(record)">编辑</a> |
        <a @click="design(record)">设计</a> |

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

import {listDef, removeDef} from "../../../api/manage";

const columns = [
  {
    dataIndex: 'name',
    slots: { title: 'customTitle' },
    scopedSlots: { customRender: 'name' },
  },
  {
    title: '路径',
    dataIndex: 'path',
  },
  {
    title: '操作',
    key: 'action',
    scopedSlots: { customRender: 'action' },
  },
];

export default {
  name: 'Mine',
  data() {
    return {
      defs: [],
      columns
    };
  },
  computed: {

  },
  created () {
    console.log('===')
    listDef().then(res => {
      console.log('listDefs', res)
      this.defs = res.data
    })
  },
  mounted () {
  },
  methods: {
    create() {
      this.$router.push({path: '/data/mine/edit/0'});
    },
    edit(record) {
      console.log(record)
      this.$router.push({path: `/data/mine/edit/${record.id}`});
    },
    design(record) {
      console.log(record)
      this.$router.push({path: `/data/mine/design/${record.id}`});
    },
    remove(record) {
      console.log(record)
      removeDef(record.id).then(json => {
        console.log('removeDef', json)
        if (json.code == 1) {
          listDef().then(res => {
            console.log('listDefs', res)
            this.defs = res.data
          })
        }
      })
    }
  }
}
</script>

<style scoped>

</style>
