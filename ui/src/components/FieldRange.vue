<template>
  <div class="panel">
    <a-form-model ref="editForm">
      <a-row :gutter="cols" class="title">
        <a-col :span="col">{{ $t('form.type') }}</a-col>
        <a-col :span="col">{{ $t('form.value') }}</a-col>
        <a-col :span="col">{{ $t('form.opt') }}</a-col>
      </a-row>

      <a-row v-if="!sections || sections.length == 0" :gutter="cols">
        <a-col :span="col"></a-col>
        <a-col :span="col"></a-col>
        <a-col :span="col">
          <a class="edit">
            <a @click="insertSection()" class="edit">{{ $t('action.add') }}</a>
          </a>
        </a-col>
      </a-row>

      <a-row v-for="item in sections" :key="item.id" :gutter="cols">

        <a-col :span="col">
          <a-form-model-item prop="type" :wrapperCol="wrapperColFull">
            <a-select v-model="item.type">
              <a-select-option value="interval">{{ $t('form.type.interval') }}</a-select-option>
              <a-select-option value="list">{{ $t('form.type.list') }}</a-select-option>
              <a-select-option value="literal">{{ $t('form.type.literal') }}</a-select-option>
            </a-select>
          </a-form-model-item>
        </a-col>

        <a-col :span="col">
          <a-form-model-item prop="value" :wrapperCol="wrapperColFull">
            <a-input v-model="item.value" />
          </a-form-model-item>
        </a-col>

        <a-col :span="8">
          <a class="edit">
            <a @click="insertSection(item)" class="edit">{{ $t('action.add') }}</a> &nbsp;
            <a @click="editSection(item)" class="edit">{{ $t('action.edit') }}</a> &nbsp;
            <a-popconfirm
                :title="$t('tips.delete')"
                :okText="$t('msg.yes')"
                :cancelText="$t('msg.no')"
                @confirm="removeSection(item)"
            >
              <a class="edit">{{ $t('action.delete') }}</a>
            </a-popconfirm>
          </a>
        </a-col>
      </a-row>
    </a-form-model>

    <a-modal
        :title="editTitle"
        :width="600"
        :visible="editSectionVisible"
        :okText="$t('form.save')"
        :cancelText="$t('form.cancel')"
        @ok="saveSection"
        @cancel="cancelSection">
      <div>
        <a-form-model ref="editForm" :model="section" :rules="rules">
          <div v-if="section.type==='interval'">
            <a-row :gutter="cols">
              <a-col :span="cols">
                <a-form-model-item :label="$t('form.start')" prop="start" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-input v-model="section.start" :placeholder="$t('tips.number.or.letter')" />
                </a-form-model-item>
              </a-col>
            </a-row>
            <a-row :gutter="cols">
              <a-col :span="cols">
                <a-form-model-item :label="$t('form.end')" prop="end" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-input v-model="section.end" :placeholder="$t('tips.number.or.letter')" />
                </a-form-model-item>
              </a-col>
            </a-row>
            <a-row :gutter="cols">
              <a-col :span="cols">
                <a-form-model-item :label="$t('form.repeat')" prop="repeat" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-input v-model="section.repeat" :precision="0" :min="1" placeholder="" />
                </a-form-model-item>
              </a-col>
            </a-row>
            <a-row :gutter="cols">
              <a-col :span="cols">
                <a-form-model-item :label="$t('form.rand')" prop="rand" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-switch v-model="section.rand" />
                </a-form-model-item>
              </a-col>
            </a-row>
            <a-row :gutter="cols" v-if="!section.rand">
              <a-col :span="cols">
                <a-form-model-item :label="$t('form.step')" prop="step" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-input v-model="section.step" :placeholder="$t('tips.number')" />
                </a-form-model-item>
              </a-col>
            </a-row>
          </div>

          <div v-if="section.type==='list'">
            <a-row :gutter="cols">
              <a-col :span="cols">
                <a-form-model-item :label="$t('form.type.list')" prop="text" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-input v-model="section.text" type="textarea" rows="3" />
                  {{$t('tips.value.each.line')}}
                </a-form-model-item>
              </a-col>
            </a-row>
          </div>

          <div v-if="section.type==='literal'">
            <a-row :gutter="cols">
              <a-col :span="cols">
                <a-form-model-item label="常量" prop="text" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
                  <a-input v-model="section.text" placeholder="" />
                </a-form-model-item>
              </a-col>
            </a-row>
          </div>
        </a-form-model>

      </div>

    </a-modal>

  </div>
</template>

<script>
import {
  listSection, createSection, removeSection, updateSection,
} from "../api/section";
import {sectionStrToArr, trimChar} from "../api/utils";

export default {
  name: 'FieldRangeComponent',
  data() {
    return {
      cols: 24,
      colsHalf: 12,
      col: 8,
      labelColFull: { lg: { span: 4 }, sm: { span: 4 } },
      labelColHalf: { lg: { span: 8}, sm: { span: 8 } },
      wrapperColFull: { lg: { span: 16 }, sm: { span: 16 } },
      wrapperColHalf: { lg: { span: 12 }, sm: { span: 12 } },

      sections: [],
      section: {},
      editTitle: '',
      editSectionVisible: false,

      rules: {
        start: [
          { required: true, message: '必须是数字或单个字母', trigger: 'change' },
          { validator: this.checkRange, trigger: 'change' },
        ],
        end: [
          { required: true, message: '必须是数字或单个字母', trigger: 'change' },
          { validator: this.checkRange, trigger: 'change' },
        ],
        repeat: [
          { validator: this.checkRepeat, message: '必须是正整数', trigger: 'change' },
        ],
        step: [
          { validator: this.checkStep, message: '必须是数字', trigger: 'change' },
        ],
      },
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
  },
  created () {
    console.log('created')

    this.loadData()
    this.$watch('time2', () => {
      console.log('time2 changed', this.time2)
      this.loadData()
    })
  },
  mounted () {
    console.log('mounted')
  },
  methods: {
    loadData () {
      console.log('loadData', this.type, this.model)
      if (!this.model.id) return

      listSection(this.model.id, this.type).then(res => {
        console.log('listSection', res)
        this.sections = res.data
      })
    },
    insertSection (item) {
      createSection(this.model.id, item ? item.id : 0, this.type).then(res => {
        console.log('createSection', res)
        this.sections = res.data
      })
    },

    editSection (item) {
      console.log('editSection', item)

      if (item.type === 'interval') {
        this.editTitle = '编辑范围'
      } else if (item.type === 'literal') {
        this.editTitle = '编辑字面常量'
        item.text = trimChar(item.value, '`')
      } else if (item.type === 'list') {
        this.editTitle = '编辑数组'
        item.text = item.value
        item.text = sectionStrToArr(item.value)
      }

      this.section = item
      this.editSectionVisible = true
    },
    saveSection() {
      this.$refs.editForm.validate(valid => {
        console.log(valid, this.section)
        if (!valid) {
          console.log('validation fail')
          return
        }

        if (this.section.type === 'interval') {
          this.section.value = this.section.start + '-' + this.section.end

          if (this.section.rand) {
            this.section.value += ':R'
            this.section.step = 0
          } else if (this.section.step && this.section.step != '' && this.section.step != 1) {
            const regx = /^[a-z,A-Z]$/
            if (regx.test(this.section.start) || regx.test(this.section.end)) {
              this.section.step = Math.floor(this.section.step)
            }

            this.section.value += ':' + this.section.step
          }

          if (this.section.repeat && this.section.repeat != '' && this.section.repeat != '1') {
            this.section.value += '{' + this.section.repeat + '}'
          }

        } else if (this.section.type === 'literal') {
          this.section.value = '`' + this.section.text + '`'

        } else if (this.section.type === 'list') {
          const arr = this.section.text.split('\n')
          this.section.value = '[' + arr.join(',') + ']'
        }

        this.section.step = parseInt(this.section.step)
        updateSection(this.section, this.type).then(res => {
          console.log('updateSection', res)
          this.sections = res.data
        })

        this.editSectionVisible = false
      })
    },
    cancelSection() {
      console.log('cancelSection')
      this.editSectionVisible = false
    },

    removeSection (item) {
      console.log(item)
      removeSection(item.id, this.type).then(res => {
        console.log('removeSection', res)
        this.sections = res.data
      })
    },
    checkRange (rule, value, callback){
      console.log('checkRange', value)

      const test1 = /^[0-9]+\.?[0-9]*$/.test(value);
      const test2 = /^[a-z,A-Z]$/.test(value);
      if (!test1 && !test2) {
        callback('必须是数字或单个字母')
      }

      callback()
    },
    checkRepeat(rule, value, callback) {
      const test = /^[1-9][0-9]*$/.test(value);
      if (!test) {
        callback('必须是正整数')
      }
      callback()
    },
    checkStep(rule, value, callback) {
      const test = /^[0-9]+\.?[0-9]*$/.test(value);
      if (!test) {
        callback('必须是数字')
      }
      callback()
    }
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
