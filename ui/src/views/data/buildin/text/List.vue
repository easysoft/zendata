<template>
  <div class="main-table">
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
        <a @click="edit(record)" :title="$t('action.edit')"><a-icon type="form" :style="{fontSize: '16px'}" /></a> &nbsp;

        <a-popconfirm
            :title="$t('tips.delete')"
            :okText="$t('msg.yes')"
            :cancelText="$t('msg.no')"
            @confirm="remove(record)"
        >
          <a href="#" :title="$t('action.delete')"><a-icon type="delete" :style="{fontSize: '16px'}" /></a> &nbsp;
        </a-popconfirm> &nbsp;

        <a-tooltip placement="top" overlayClassName="tooltip-light">
          <template slot="title">
            <div class="content-width" style="min-width: 280px;">
              <div class="title">{{$t('tips.refer')}}</div>
              <div class="content">
                <div>range: {{ record.referName }}</div>
              </div>
            </div>
          </template>
          <a href="#" :title="$t('tips.refer')">&nbsp; <a-icon type="link" :style="{fontSize: '16px'}" /></a>
        </a-tooltip>

      </span>

    </a-table>

    <div class="pagination-wrapper">
      <a-pagination @change="onPageChange" :current="page" :total="total" :defaultPageSize="15" simple size="small" />
    </div>

    <a-modal
      :visible="editModalVisible"
      :title="editModalVisible ? editRecord ? `${$t('menu.text.edit')}: ${editRecord.title}` : $t('title.text.create') : ''"
      :footer="false"
      :centered="true"
      :width="700"
      @cancel="handleCancelEditModal"
    >
      <Edit
        :v-if="editModalVisible"
        :id="editModalVisible ? editID ? editID : 0 : null"
        :afterSave="handleEditSave"
      />
    </a-modal>
  </div>
</template>

<script>

import {listText, removeText} from "../../../../api/manage";
import {PageSize, pathToRelated} from "../../../../api/utils";
import debounce from "lodash.debounce"
import Edit from './Edit';

export default {
  name: 'TextList',
  components: {
    Edit,
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
        width: 450
      },
      {
        title: this.$i18n.t('form.opt'),
        key: 'action',
        scopedSlots: { customRender: 'action' },
        width: 80
      },
    ];

    return {
      models: [],
      columns,

      designVisible: false,
      designModel: {},
      time: 0,

      page: 1,
      total: 0,
      pageSize: PageSize,
    };
  },
  computed: {
    keywords: function() {
      if (this.$route.query && typeof this.$route.query.search === 'string') {
        return this.$route.query.search;
      }
      return '';
    },
    editID: function() {
      if (this.$route.params && this.$route.params.id !== undefined) {
        return this.$route.params.id;
      }
      return null;
    },
    editModalVisible: function() {
      return this.editID !== null;
    },
    editRecord: function() {
      const {editID} = this;
      if (!editID) {
        return null;
      }
      return this.models.find(x => x.id == editID);
    }
  },
  watch: {
    keywords: function() {
      this.onSearch();
    }
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
      this.$router.push({path: '/data/buildin/text/edit/0'});
    },
    loadData() {
      listText(this.keywords, this.page).then(json => {
        console.log('listText', json)
        this.models = json.data
        this.total = json.total
      })
    },
    handleCancelEditModal() {
      const {path, query} = this.$route;
      const newPath = '/data/buildin/text/list';
      if (path !== newPath) {
        this.$router.replace({path: newPath, query});
      }
    },
    handleEditSave() {
      this.handleCancelEditModal();
      this.loadData();
    },
    edit(record) {
      const {path, query = {}} = this.$router;
      const newPath = `/data/buildin/text/list/${record.id}`;
      if (path !== newPath) {
        this.$router.replace({path: newPath, query});
      }
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
