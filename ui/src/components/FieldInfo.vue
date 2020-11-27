<template>
  <div>
    <a-form-model ref="editForm" :model="model" :rules="rules">
      <a-row :gutter="colsFull">
          <a-form-model-item label="名称" prop="field" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.field" />
          </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item label="前缀" :labelCol="labelColFull" class="zui-input-group">
          <a-form-model-item prop="prefix" :style="{ display: 'inline-block', width: 'calc(40% - 35px)' }">
            <a-input v-model="model.prefix" />
          </a-form-model-item>

          <span class="zui-input-group-addon" :style="{ width: '60px' }">
            <span>后缀</span>
          </span>

          <a-form-model-item prop="postfix" :style="{ display: 'inline-block', width: 'calc(27% - 30px)' }">
            <a-input v-model="model.postfix" />
          </a-form-model-item>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item label="循环" :labelCol="labelColFull" class="zui-input-group">
          <a-form-model-item prop="loop" :style="{ display: 'inline-block', width: 'calc(40% - 35px)' }">
            <a-input v-model="model.loop" placeholder="数字3，或区间1-3" />
          </a-form-model-item>

          <span class="zui-input-group-addon" :style="{ width: '60px' }">
              <span>间隔</span>
            </span>

          <a-form-model-item prop="loopfix" :style="{ display: 'inline-block', width: 'calc(27% - 30px)' }">
            <a-input v-model="model.loopfix" />
          </a-form-model-item>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item label="类型" :labelCol="labelColFull" class="zui-input-group">
          <a-form-model-item prop="type" :style="{ display: 'inline-block', width: 'calc(40% - 35px)' }">
            <a-select v-model="model.type">
              <a-select-option value="list">列表</a-select-option>
              <a-select-option value="timestamp">时间戳</a-select-option>
            </a-select>
          </a-form-model-item>

          <span class="zui-input-group-addon" :style="{ width: '60px' }">
              <span>模式</span>
          </span>

          <a-form-model-item prop="mode" :style="{ display: 'inline-block', width: 'calc(27% - 30px)' }">
            <a-select v-model="model.mode">
              <a-select-option value="parallel">平行</a-select-option>
              <a-select-option value="recursive">递归</a-select-option>
            </a-select>
          </a-form-model-item>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item label="宽度" :labelCol="labelColFull" class="zui-input-group">
          <a-form-model-item prop="length" :style="{ display: 'inline-block', width: 'calc(40% - 45px)' }">
            <a-input v-model="model.length" :min="0" />
          </a-form-model-item>

          <span class="zui-input-group-addon" :style="{ width: '70px' }">
            <span>左占位符</span>
          </span>
          <a-form-model-item prop="loopfix" :style="{ display: 'inline-block', width: 'calc(13% - 45px)' }">
            <a-input v-model="model.leftPad" />
          </a-form-model-item>

          <span class="zui-input-group-addon" :style="{ width: '70px' }">
            <span>右占位符</span>
          </span>
          <a-form-model-item prop="loopfix" :style="{ display: 'inline-block', width: 'calc(13% - 43px)' }">
            <a-input v-model="model.rightPad" />
          </a-form-model-item>
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="格式" prop="format" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <div class="inline">
              <a-input v-model="model.format">
                <a-select v-model="model.format" slot="addonAfter" default-value="" style="width: 100px">
                  <a-select-option value="">
                    选择函数
                  </a-select-option>
                  <a-select-option value="md5">md5</a-select-option>
                  <a-select-option value="sha1">sha1</a-select-option>
                  <a-select-option value="base64">base64</a-select-option>
                  <a-select-option value="urlencode">urlencode</a-select-option>
                </a-select>
              </a-input>
            </div>
          </a-form-model-item>

        </a-col>
        <a-col :span="colsHalf">
          <a-form-model-item label="右占位符" prop="rand" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
            <a-switch v-model="model.rand" />
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row :gutter="colsFull">
          <a-form-model-item label="描述" prop="note" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.note" type="textarea" rows="3" />
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
import {saveDefField, saveInstancesItem} from "../api/manage";
import {checkLoop} from "../api/utils";
import {labelColLarge, wrapperColLarge} from "../utils/const";

export default {
  name: 'FieldInfoComponent',
  data() {
    return {
      labelCol: labelColLarge,
      wrapperCol: wrapperColLarge,

      colsFull: 24,
      colsHalf: 12,
      labelColFull: { lg: { span: 4 }, sm: { span: 4 } },
      wrapperColFull: { lg: { span: 16 }, sm: { span: 16 } },
      labelColHalf: { lg: { span: 8}, sm: { span: 8 } },
      labelColHalf2: { lg: { span: 4}, sm: { span: 4 } },
      wrapperColHalf: { lg: { span: 12 }, sm: { span: 12 } },
      rules: {
        field: [
          { required: true, message: '名称不能为空', trigger: 'change' },
        ],
        loop: [
          { validator: checkLoop, message: '需为正整数或其区间', trigger: 'change' },
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
  },

  computed: {
  },
  created () {
    console.log('created')
  },
  mounted () {
    console.log('mounted')
  },
  methods: {
    save() {
      console.log('save')
      this.$refs.editForm.validate(valid => {
        console.log(valid, this.model)
        if (!valid) {
          console.log('validation fail')
          return
        }

        if (this.type === 'def') {
          saveDefField(this.model).then(json => {
            console.log('saveDefField', json)
            this.$emit('save')
          })
        } else {
          saveInstancesItem(this.model).then(json => {
            console.log('saveInstancesItem', json)
            this.$emit('save')
          })
        }
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
</style>
