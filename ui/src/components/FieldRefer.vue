<template>
  <div class="panel">
    <a-form-model ref="editForm" :model="refer" :rules="rules">
      <a-row :gutter="colsFull">
        <a-form-model-item label="类型" prop="type" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.type" @change="onTypeChanged">
            <a-select-option value="ranges">Ranges</a-select-option>
            <a-select-option value="instances">Instances</a-select-option>
            <a-select-option value="config">Config</a-select-option>
            <a-select-option value="yaml">Yaml</a-select-option>

            <a-select-option value="excel">Excel</a-select-option>
            <a-select-option value="text">Text</a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="文件" prop="file" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.file">
              <a-select slot="addonAfter" default-value="" style="width: 80px">
                <a-select-option value="">
                  选择
                </a-select-option>
                <a-select-option v-for="f in files" :key="f.path">
                  {{ f.name }}
                </a-select-option>
              </a-select>
            </a-input>
          </a-form-model-item>
        </a-col>

        <a-col :span="colsHalf">
          <a-form-model-item v-if="!showColIndex" label="列名" prop="colName" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.colName">
              <a-select slot="addonAfter" default-value="" style="width: 80px">
                <a-select-option value="">
                  选择
                </a-select-option>
              </a-select>
            </a-input>
          </a-form-model-item>
          <a-form-model-item v-if="showColIndex" label="列索引" prop="colIndex" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.colIndex">
              <a-select slot="addonAfter" default-value="" style="width: 80px">
                <a-select-option value="">
                  选择
                </a-select-option>
              </a-select>
            </a-input>
          </a-form-model-item>
        </a-col>
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
import {getDefFieldRefer, updateDefFieldRefer} from "../api/manage";

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
        start: [
          { required: true, message: '必须是数字或单个字母', trigger: 'change' },
          { validator: this.checkRange, trigger: 'change' },
        ],
      },

      res: {},
      files: [],
      fields: [],
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
    showColIndex() {
      return this.refer.type == 'text' && !this.refer.hasTitle
    }
  },
  created () {
    console.log('created')

    this.loadDefFieldRefer()
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadDefFieldRefer()
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
        this.res = json.res
      })
    },
    onTypeChanged() {
      console.log('onTypeChanged')
      this.files = this.res[this.refer.type]
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
  }
}
</script>

<style lang="less" scoped>
  .panel {
    padding: 4px 8px;

  }
</style>
