<template>
  <div>
    <div class="head">
      <div class="title">{{ $t('menu.instances.list') }}</div>
      <div class="filter">
        <a-input-search v-model="keywords" @change="onSearch" :allowClear="true"
                        placeholder="$t('tips.search')" style="width: 300px" />
      </div>
      <div class="buttons">
        <a-button type="primary" @click="create()">{{ $t('action.create') }}</a-button>
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
        <a @click="edit(record)">{{ $t('action.edit') }}</a> |
        <a @click="design(record)">{{ $t('action.design') }}</a> |

        <a-popconfirm
            :title="$t('tips.delete')"
            :okText="$t('msg.yes')"
            :cancelText="$t('msg.no')"
            @confirm="remove(record)"
        >
          <a href="#">{{ $t('action.delete') }}</a>
        </a-popconfirm> |

        <a-tooltip placement="top" overlayClassName="tooltip-light">
          <template slot="title">
            <div class="content-width" style="min-width: 280px;">
              <div class="title">{{$t('tips.refer')}}</div>
              <div class="content">
                <div>from: {{ record.referName }}</div>
                <div>use: field_name</div>
              </div>
            </div>
          </template>
          <a href="#">{{$t('tips.refer')}}</a>
        </a-tooltip>

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

import {listInstances, removeInstances} from "../../../../api/manage";
import { DesignComponent } from '../../../../components'
import {PageSize, ResTypeInstances, pathToRelated} from "../../../../api/utils";
import debounce from "lodash.debounce"

export default {
  name: 'InstanceList',
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
      models: [],
      columns,

      designVisible: false,
      designModel: {},
      type: ResTypeInstances,
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
      this.$router.push({path: '/data/buildin/instances/edit/0'});
    },
    loadData() {
      listInstances(this.keywords, this.page).then(json => {
        console.log('listInstances', json)
        this.models = json.data
        this.total = json.total
      })
    },
    edit(record) {
      console.log(record)
      this.$router.push({path: `/data/buildin/instances/edit/${record.id}`});
    },
    design(record) {
      this.time = Date.now() // trigger data refresh
      console.log(record)
      this.designVisible = true
      this.designModel = record
    },
    remove(record) {
      console.log(record)
      removeInstances(record.id).then(json => {
        console.log('removeInstances', json)
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
