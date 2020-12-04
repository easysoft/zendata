<template>
  <div class="panel">
    <a-form-model ref="editForm" :model="refer" :rules="rules">
      <a-row :gutter="colsFull">
        <a-form-model-item label="类型" prop="type" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.type" @change="onReferTypeChanged">
            <a-select-option value="config">字段</a-select-option>
            <a-select-option value="ranges">序列</a-select-option>
            <a-select-option value="instances">实例</a-select-option>
            <a-select-option value="yaml">执行</a-select-option>
            <a-select-option value="excel">表格</a-select-option>
            <a-select-option value="text">文本</a-select-option>
            <a-select-option value="value">表达式</a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row v-if="refer.type && refer.type!='value'" :gutter="colsFull">
          <a-form-model-item label="文件" prop="file" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-select v-model="refer.file" @change="onReferFileChanged">
              <a-select-option value="">选择</a-select-option>
              <a-select-option v-for="(f, i) in files" :value="f.referName" :key="i">
                <span v-if="refer.type != 'excel'">{{ f.title }}</span>
                <span v-if="refer.type == 'excel'">{{ f.referName }}</span>
              </a-select-option>
            </a-select>
          </a-form-model-item>
      </a-row>

      <a-row v-if="refer.type==='excel'" :gutter="colsFull">
        <a-form-model-item label="Excel表格" prop="sheet" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.sheet" @change="onReferSheetChanged">
            <a-select-option value="">选择</a-select-option>
            <a-select-option v-for="(f, i) in sheets" :value="f.sheet" :key="i">
              {{ f.sheet }}
            </a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row v-if="showColSection" :gutter="colsFull">
        <a-form-model-item label="列名" prop="colName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
           <a-select v-model="refer.colName">
              <a-select-option value="">选择</a-select-option>
              <a-select-option v-for="f in fields" :key="f.name">
                {{ f.name }}
              </a-select-option>
            </a-select>
        </a-form-model-item>
      </a-row>

      <a-row v-if="showCount" :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="取记录数" prop="count" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.count" />
            0表示取所有记录
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row v-if="showStep" :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="步长" prop="step" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.step" :disabled="refer.rand" />
          </a-form-model-item>
        </a-col>
        <a-col :span="colsHalf">
          <a-form-model-item label="是否随机" prop="rand" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-switch v-model="refer.rand" />
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row v-if="refer.type=='value'" :gutter="colsFull">
        <a-form-model-item label="表达式" prop="value" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-input v-model="refer.value" />
          <span class="input-tips">
            请输入数学运算表达式，由相同文件中的字段组成，如"($field_step_negative * $field_nested_range) * -1 + 1000"
          </span>
        </a-form-model-item>
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
import {listReferFileForSelection, listReferSheetForSelection,
  listReferResFieldForSelection, listReferExcelColForSelection, getRefer, updateRefer,
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
    showStep() {
      return this.refer.type === 'text'
    },
    showCount() {
      return this.refer.type === 'yaml' || this.refer.type === 'ranges' || this.refer.type === 'instances'
    },
    showColSection() {
      return this.refer.type === 'ranges' || this.refer.type === 'instances' || this.refer.type === 'excel'
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
      if (!this.model.id) return

      getRefer(this.model.id, this.type).then(json => {
        console.log('getRefer', json)
        this.refer = json.data

        this.removeSheet()
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
        this.listReferSheetForSelection(false)
      } else {
        this.listReferResFieldForSelection(false)
      }
    },
    onReferSheetChanged() {
      console.log("onReferSheetChanged")

      if (this.refer.type == 'excel') {
        this.listReferResFieldForSelection(false)
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

        let data = JSON.parse(JSON.stringify(this.refer))
        if (data.type === 'excel') data.file = data.file + '.' + data.sheet

        data.count = parseInt(data.count)
        data.step = parseInt(data.step)
        updateRefer(data, this.type).then(json => {
          console.log('updateRefer', json)
        })
      })
    },
    reset() {
      console.log('reset')
      this.$refs.editForm.reset()
    },

    listReferFileForSelection(resType, init) {
      if (!this.refer.type) return

      listReferFileForSelection(resType).then(json => {
        console.log('listReferFileForSelection', json)
        this.files = json.data

        if (init) {
          if (this.refer.type === 'excel') {
            this.listReferSheetForSelection(init)
          } else {
            this.listReferResFieldForSelection(init)
          }
        } else {
          this.refer.file = ''
        }
      })
    },
    listReferSheetForSelection(init) {
      if (!this.refer.type) return

      listReferSheetForSelection(this.refer.file + '.' + this.refer.sheet ).then(json => {
        console.log('listReferSheetForSelection', json)
        this.sheets = json.data

        if (init) {
          this.listReferResFieldForSelection(true)
        } else {
          this.refer.sheet = ''
        }
      })
    },
    listReferResFieldForSelection(init) {
      if (!this.refer.type) return

      let id = 0
      if (this.refer.type === 'excel') {
        listReferExcelColForSelection(this.refer.file + '.' + this.refer.sheet).then(json => {
          console.log('listReferExcelColForSelection', json)
          this.fields = json.data

          if (!init) {
            this.refer.colName = ''
          }
        })

      } else if (this.refer.type != 'value') {
        this.files.forEach((fi) => {
          if (fi.referName === this.refer.file) id = fi.id
        })

        listReferResFieldForSelection(id, this.refer.type).then(json => {
          console.log('listReferResFieldForSelection', json)
          this.fields = json.data

          if (!init) {
            this.refer.colName = ''
          }
        })
      }
    },
    removeSheet() {
      if (this.refer.type == 'excel') {
        this.refer.file = this.refer.file.substring(0, this.refer.file.lastIndexOf('.'))
      }
    }
  }
}
</script>

<style lang="less" scoped>
  .panel {
    padding: 4px 8px;
  }
</style>
