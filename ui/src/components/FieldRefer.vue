<template>
  <div class="panel">
    <a-form-model ref="editForm" :model="refer" :rules="rules">
      <a-row :gutter="colsFull">
        <a-form-model-item label="类型" prop="type" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.type" @change="onTypeChanged">
            <a-select-option value="ranges">序列（Ranges）</a-select-option>
            <a-select-option value="instances">实例（Instances）</a-select-option>
            <a-select-option value="config">配置（Config）</a-select-option>
            <a-select-option value="yaml">内容（来自YAML）</a-select-option>
            <a-select-option value="excel">表格（Excel）</a-select-option>
            <a-select-option value="text">文本（Text）</a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
          <a-form-model-item label="文件" prop="file" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="refer.file">
              <a-select v-model="referFile" @change="onReferChanged" slot="addonAfter" style="width: 300px">
                <a-select-option value="">选择</a-select-option>
                <a-select-option v-for="(f, i) in files" :value="f.name+'-'+f.title" :key="i">
                  {{ f.title }}
                </a-select-option>
              </a-select>
            </a-input>
          </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item v-if="!showColSection" label="列名" prop="colName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-input v-model="refer.colName">
            <a-select v-model="referFieldName" @change="onFieldNameChanged" slot="addonAfter" style="width: 300px">
              <a-select-option value="">选择</a-select-option>
              <a-select-option v-for="f in fields" :key="f.name">
                {{ f.name }}
              </a-select-option>
            </a-select>
          </a-input>
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
import {getDefFieldRefer, updateDefFieldRefer, listDefFieldReferType, listDefFieldReferField} from "../api/manage";

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
      fields: [],
      referFile: '',

      referFieldName: '',
      referFieldIndex: '',
    };
  },
  props: {
    field: {
      type: Object,
      default: () => null
    },
    time: {
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

    this.loadDefFieldRefer()
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadDefFieldRefer("")
    })
  },
  mounted () {
    console.log('mounted')
  },
  methods: {
    loadDefFieldRefer() {
      getDefFieldRefer(this.field.id).then(json => {
        console.log('getDefFieldRefer', json)
        this.refer = json.data
        this.listDefFieldReferType(this.refer.type, true)
      })
    },

    onTypeChanged() {
      console.log('onTypeChanged')
      this.listDefFieldReferType(this.refer.type, false)
    },
    onReferChanged(value) {
      console.log("onReferChanged")
      this.refer.file = value.split('-')[0]

      if (this.refer.type != 'yaml' && this.refer.type != 'text') {
        this.listDefFieldReferField()
      } else {
        this.refer.colName = ''
      }
    },
    onFieldNameChanged(value) {
      console.log("onFieldChanged")
      this.refer.colName = value
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
        updateDefFieldRefer(this.refer).then(json => {
          console.log('updateDefFieldRefer', json)
        })
      })
    },
    reset() {
      console.log('reset')
      this.$refs.editForm.reset()
    },

    listDefFieldReferType(resType, init) {
      listDefFieldReferType(resType).then(json => {
        console.log('listDefFieldReferType', json)
        this.files = json.data
      })

      if (!init) {
        this.refer.file = ''
        this.referFile = ''
      }
    },
    listDefFieldReferField() {
      let file = {}
      for (let i = 0; i < this.files.length; i++) {
        const f = this.files[i]
        if (f.name + '-' + f.title === this.referFile) {
          file = f
          break
        }
      }

      listDefFieldReferField(file).then(json => {
        console.log('listDefFieldReferField', json)
        this.fields = json.data
      })
      this.refer.colName = ''
      this.refer.colIndex = ''

      this.referFieldName = ''
      this.referFieldIndex = ''
    }
  }
}
</script>

<style lang="less" scoped>
  .panel {
    padding: 4px 8px;

  }
</style>
