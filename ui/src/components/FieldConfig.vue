<template>
  <div class="panel">
    <div class="radios">
      <a-radio-group :value="model.isRange" button-style="solid">
        <a-radio-button :value="true">
          区间
        </a-radio-button>
        <a-radio-button :value="false">
          引用
        </a-radio-button>
      </a-radio-group>
      &nbsp;&nbsp;&nbsp;
      <span class="range">{{model.range}}</span>
    </div>
    <div>
      <a-form-model ref="editForm">
        <a-row :gutter="cols" class="title">
          <a-col :span="col">取值</a-col>
          <a-col :span="col">类型</a-col>
          <a-col :span="col">操作</a-col>
        </a-row>

        <a-row v-if="!sections || sections.length == 0" :gutter="cols">
          <a-col :span="col"></a-col>
          <a-col :span="col"></a-col>
          <a-col :span="8">
            <a class="edit">
              <a @click="insertSection()" class="edit">添加</a>
            </a>
          </a-col>
        </a-row>

        <a-row v-for="item in sections" :key="item.id" :gutter="cols">
          <a-col :span="col">
            <a-form-model-item prop="range" :wrapperCol="wrapperColFull">
              <a-input v-model="item.value" />
            </a-form-model-item>
          </a-col>
          <a-col :span="col">
            <a-form-model-item prop="range" :wrapperCol="wrapperColFull">
              <a-select v-model="item.type">
                <a-select-option value="scope">范围</a-select-option>
                <a-select-option value="arr">数组</a-select-option>
                <a-select-option value="const">字面常量</a-select-option>
              </a-select>
            </a-form-model-item>
          </a-col>
          <a-col :span="8">
            <a class="edit">
              <a @click="insertSection(item)" class="edit">添加</a> |
              <a @click="editSection(item)" class="edit">编辑</a> |
              <a @click="removeSection(item)" class="edit">删除</a>
            </a>
          </a-col>
        </a-row>
      </a-form-model>
    </div>

    <a-modal
        :title="editTitle"
        :width="600"
        :visible="editSectionVisible"
        okText="保存"
        cancelText="取消"
        @ok="saveSection"
        @cancel="cancelSection">
      <div>
        <div v-if="section.type==='scope'">
          <a-row :gutter="cols">
            <a-col :span="cols">
              <a-form-model-item label="开始" prop="prefix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-input v-model="section.start" placeholder="数字或单个字母" />
              </a-form-model-item>
            </a-col>
          </a-row>
          <a-row :gutter="cols">
            <a-col :span="cols">
              <a-form-model-item label="结束" prop="postfix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-input v-model="section.end" placeholder="数字或单个字母" />
              </a-form-model-item>
            </a-col>
          </a-row>
          <a-row :gutter="cols">
            <a-col :span="cols">
              <a-form-model-item label="重复次数" prop="prefix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-input v-model="section.repeat" :min="1" placeholder="" />
              </a-form-model-item>
            </a-col>
          </a-row>
          <a-row :gutter="cols">
            <a-col :span="cols">
              <a-form-model-item label="随机" prop="prefix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-switch v-model="section.rand" />
              </a-form-model-item>
            </a-col>
          </a-row>
          <a-row :gutter="cols" v-if="!section.rand">
            <a-col :span="cols">
              <a-form-model-item label="步长" prop="postfix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-input v-model="section.step" placeholder="数字" />
              </a-form-model-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="section.type==='arr'">
          <a-row :gutter="cols">
            <a-col :span="cols">
              <a-form-model-item label="数组" prop="prefix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-input v-model="section.text" type="textarea" rows="3" />
                每行一个值
              </a-form-model-item>
            </a-col>
          </a-row>
        </div>

        <div v-if="section.type==='const'">
          <a-row :gutter="cols">
            <a-col :span="cols">
              <a-form-model-item label="字面常量" prop="prefix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                <a-input v-model="section.text" placeholder="" />
              </a-form-model-item>
            </a-col>
          </a-row>
        </div>
      </div>
    </a-modal>

  </div>
</template>

<script>
import { listDefFieldSection, createDefFieldSection, updateDefFieldSection, removeDefFieldSection } from "../api/manage";

export default {
  name: 'FieldConfigComponent',
  data() {
    return {
      cols: 24,
      colsHalf: 12,
      col: 8,
      labelColFull: { lg: { span: 4 }, sm: { span: 4 } },
      labelColHalf: { lg: { span: 8}, sm: { span: 8 } },
      wrapperColFull: { lg: { span: 18 }, sm: { span: 18 } },
      wrapperColHalf: { lg: { span: 12 }, sm: { span: 12 } },

      sections: [],
      section: {},
      editTitle: '',
      editSectionVisible: false,
    };
  },
  props: {
    model: {
      type: Object,
      default: () => null
    },
    time: {
      type: Number,
      default: () => 0
    },
  },

  computed: {
  },
  created () {
    console.log('created')
    this.loadData()
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadData()
    })
  },
  mounted () {
    console.log('mounted1')
  },
  methods: {
    loadData () {
      if (!this.model.id) return

      listDefFieldSection(this.model.id).then(res => {
        console.log('listDefFieldSection', res)
        this.sections = res.data
      })
    },
    insertSection (item) {
      console.log(item)
      createDefFieldSection(this.model.id, item?item.id:0).then(res => {
        console.log('createDefFieldSection', res)
        this.sections = res.data
      })
    },

    editSection (item) {
      console.log(item)
      if (item.type === 'scope') {
        this.editTitle = '编辑范围'
      } else if (item.type === 'arr') {
        this.editTitle = '编辑数组'
      } else if (item.type === 'const') {
        this.editTitle = '编辑字面常量'
      }

      this.section = item
      this.editSectionVisible = true
    },
    saveSection() {
      console.log('saveSection', this.section)

      if (this.section.type === 'scope') {
        this.section.value = this.section.start + '-' + this.section.end

        if (this.section.rand) {
          this.section.value += ':R'
          this.section.step = 0
        } else if (this.section.step && this.section.step != '' && this.section.step != 1) {
          this.section.value += ':' + this.section.step
        }

        if (this.section.repeat && this.section.repeat != '' && this.section.repeat != '1') {
          this.section.value += '{' + this.section.repeat + '}'
        }

      } else if (this.section.type === 'arr') {
        const arr = this.section.text.split('\n')
        this.section.value = '[' + arr.join(',') + ']'

      } else if (this.section.type === 'const') {
        this.section.value = '`' + this.section.text + '`'
      }

      updateDefFieldSection(this.section).then(res => {
        console.log('updateDefFieldSection', res)
        this.sections = res.data
      })

      this.editSectionVisible = false
    },
    cancelSection() {
      console.log('cancelSection')
      this.editSectionVisible = false
    },

    removeSection (item) {
      console.log(item)

      removeDefFieldSection(item.id).then(res => {
        console.log('removeDefFieldSection', res)
        this.sections = res.data
      })
    },
  }
}
</script>

<style lang="less" scoped>
  .panel {
    padding: 4px 8px;
    .title {
      font-weight: bolder;
      margin-bottom: 5px;
      padding-bottom: 5px;
      border-bottom: 1px solid #e9f2fb;
    }
    .radios {
      margin-bottom: 12px;
      .range {
        display: inline-block;
        margin-left: 12px;
      }
    }
    .edit {
      line-height: 32px;
    }
  }
</style>
