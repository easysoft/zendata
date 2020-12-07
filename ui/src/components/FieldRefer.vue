<template>
  <div class="panel">
    <a-form-model ref="editForm" :model="refer" :rules="rules">
      <a-row :gutter="colsFull">
        <a-form-model-item :label="$t('form.type')" prop="type" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.type" @change="onReferTypeChanged">
            <a-select-option value="config">{{$t('msg.config')}}</a-select-option>
            <a-select-option value="ranges">{{$t('msg.ranges')}}</a-select-option>
            <a-select-option value="instances">{{$t('msg.instances')}}</a-select-option>
            <a-select-option value="yaml">{{$t('msg.exec')}}</a-select-option>
            <a-select-option value="excel">{{$t('msg.excel')}}</a-select-option>
            <a-select-option value="text">{{$t('msg.text')}}</a-select-option>
            <a-select-option value="value">{{$t('form.expr')}}</a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row v-if="refer.type && refer.type!='value'" :gutter="colsFull">
          <a-form-model-item :label="$t('form.file')" prop="file" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-select v-model="refer.file" @change="onReferFileChanged">
              <a-select-option value="">{{$t('tips.pls.select')}}</a-select-option>
              <a-select-option v-for="(f, i) in files" :value="f.referName" :key="i">
                <span v-if="refer.type != 'excel'">{{ f.title }}</span>
                <span v-if="refer.type == 'excel'">{{ f.referName }}</span>
              </a-select-option>
            </a-select>
          </a-form-model-item>
      </a-row>

      <a-row v-if="refer.type==='excel'" :gutter="colsFull">
        <a-form-model-item :label="$t('msg.excel.sheet')" prop="sheet" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-select v-model="refer.sheet" @change="onReferSheetChanged">
            <a-select-option value="">{{$t('tips.pls.select')}}</a-select-option>
            <a-select-option v-for="(f, i) in sheets" :value="f.sheet" :key="i">
              {{ f.sheet }}
            </a-select-option>
          </a-select>
        </a-form-model-item>
      </a-row>

      <a-row v-if="showColSection" :gutter="colsFull">
        <a-form-model-item :label="$t('form.col')" prop="colName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
           <a-select v-model="refer.colName">
              <a-select-option value="">{{$t('tips.pls.select')}}</a-select-option>
              <a-select-option v-for="f in fields" :key="f.name">
                {{ f.name }}
              </a-select-option>
            </a-select>
        </a-form-model-item>
      </a-row>

      <a-row v-if="showCount" :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item :label="$t('form.count')" prop="count" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.count" />
            {{$t('tips.zero')}}
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row v-if="showStep" :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item :label="$t('form.step')" prop="step" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="refer.step" :disabled="refer.rand" />
          </a-form-model-item>
        </a-col>
        <a-col :span="colsHalf">
          <a-form-model-item :label="$t('form.rand')" prop="rand" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-switch v-model="refer.rand" />
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row v-if="refer.type=='value'" :gutter="colsFull">
        <a-form-model-item :label="$t('form.expr')" prop="value" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-input v-model="refer.value" />
          <span class="input-tips">{{$t('tips.expr')}}</span>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item class="center">
          <a-button @click="save" type="primary">{{$t('form.save')}}</a-button>
          <a-button @click="reset" style="margin-left: 10px;">{{$t('form.reset')}}</a-button>
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
          this.refer.sheet = ''
          this.refer.colName = ''
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
          this.refer.colName = ''
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
