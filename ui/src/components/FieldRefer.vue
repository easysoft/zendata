<template>
  <div class="panel">
    <a-form-model ref="editForm" :model="refer" :rules="rules">
      <a-row :gutter="colsFull">
        <a-form-model-item label="类型" prop="type" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.type" @change="onReferTypeChanged">
            <a-select-option value="config">字段</a-select-option>
            <a-select-option value="ranges">序列</a-select-option>
            <a-select-option value="instances">实例</a-select-option>
            <a-select-option value="yaml">内容</a-select-option>
            <a-select-option value="excel">表格</a-select-option>
            <a-select-option value="text">文本</a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
          <a-form-model-item label="文件" prop="file" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-select v-model="refer.file" @change="onReferFileChanged">
              <a-select-option value="">选择</a-select-option>
              <a-select-option v-for="(f, i) in files" :value="f.id" :key="i">
                {{ f.title }}
              </a-select-option>
            </a-select>
          </a-form-model-item>
      </a-row>

      <a-row v-if="refer.type==='excel'" :gutter="colsFull">
        <a-form-model-item label="Excel表格" prop="sheet" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.sheet" @change="onReferSheetChanged">
            <a-select-option value="">选择</a-select-option>
            <a-select-option v-for="(f, i) in sheets" :value="f.id" :key="i">
              {{ f.title }}
            </a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item v-if="!showColSection" label="列名" prop="colName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
           <a-select v-model="referFieldName">
              <a-select-option value="">选择</a-select-option>
              <a-select-option v-for="f in fields" :key="f.name">
                {{ f.name }}
              </a-select-option>
            </a-select>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="取记录数" prop="count" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.count" />
            0表示取所有记录
          </a-form-model-item>
        </a-col>
        <a-col :span="colsHalf" v-if="refer.type == 'text'">
          <a-form-model-item label="是否含标题" prop="hasTitle" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
            <a-switch v-model="refer.hasTitle" />
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item class="center">
          <a-button @click="save" type="primary">保存</a-button>
          <a-button @click="reset" style="margin-left: 10px;">重置</a-button>
        </a-form-model-item>
      </a-row>

    </a-form-model>
  </div>
</template>

<script>
import {listReferFileForSelection, listReferSheetForSelection, listReferFieldForSelection, getRefer, updateRefer,
} from "../api/refer";

export default {
  name: 'FieldReferComponent',
  data() {
    return {
      colsFull: 24,
      colsHalf: 12,
      labelColFull: { lg: { span: 4 }, sm: { span: 4 } },
      wrapperColFull: { lg: { span: 16 }, sm: { span: 16 } },
      labelColHalf: { lg: { span: 8}, sm: { span: 8 } },
      labelColHalf2: { lg: { span: 4}, sm: { span: 4 } },
      wrapperColHalf: { lg: { span: 12 }, sm: { span: 12 } },

      refer: {},
      rules: {
      },

      res: {},
      files: [],
      sheets: [],
      fields: [],

      referFieldName: '',
      referFieldIndex: '',
    };
  },
  props: {
    type: {
      type: String,
      default: () => ''
    },
    model: {
      type: Object,
      default: () => null
    },
    time2: {
      type: Number,
      default: () => 0
    },
  },

  computed: {
    showColSection() {
      return this.refer.type == 'yaml' || this.refer.type == 'text'
    }
  },
  created () {
    console.log('created')

    this.loadData()
    this.$watch('time2', () => {
      console.log('time2 changed', this.time2)
      this.loadData("")
    })
  },
  mounted () {
    console.log('mounted')
  },
  methods: {
    loadData() {
      getRefer(this.model.id, this.type).then(json => {
        console.log('getRefer', json)
        this.refer = json.data
        this.listReferFileForSelection(this.refer.type, true)
      })
    },

    onReferTypeChanged() {
      console.log('onReferTypeChanged')
      this.listReferFileForSelection(this.refer.type, false)
    },
    onReferFileChanged() {
      console.log("onReferFileChanged")

      if (this.refer.type == 'excel') {
        this.listReferSheetForSelection()
      } else {
        this.listReferFieldForSelection()
      }
    },
    onReferSheetChanged() {
      console.log("onReferSheetChanged")

      if (this.refer.type == 'excel') {
        this.listReferFieldForSelection()
      }
    },

    save() {
      console.log('save')
      this.$refs.editForm.validate(valid => {
        console.log(valid, this.refer)
        if (!valid) {
          console.log('validation fail')
          return
        }

        this.refer.count = parseInt(this.refer.count)
        updateRefer(this.refer, this.type).then(json => {
          console.log('updateRefer', json)
        })
      })
    },
    reset() {
      console.log('reset')
      this.$refs.editForm.reset()
    },

    listReferFileForSelection(resType, init) {
      listReferFileForSelection(resType).then(json => {
        console.log('listReferFileForSelection', json)
        this.files = json.data
      })

      if (!init) {
        this.refer.file = ''
      }
    },
    listReferSheetForSelection(resType, init) {
      listReferSheetForSelection(resType).then(json => {
        console.log('listReferSheetForSelection', json)
        this.files = json.data
      })

      if (!init) {
        this.refer.sheet = ''
      }
    },
    listReferFieldForSelection() {
      let referTo = ''
      if (this.refer.type != 'excel') {
        referTo = this.refer.file
      } else {
        referTo = this.refer.sheet
      }
      listReferFieldForSelection(referTo, this.refer.type).then(json => {
        console.log('listReferFieldForSelection', json)
        this.fields = json.data
      })
      this.refer.colName = ''
      this.refer.colIndex = ''
    }
  }
}
</script>

<style lang="less" scoped>
  .panel {
    padding: 4px 8px;
  }
</style>
