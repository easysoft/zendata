<template>
  <div>
    <a-form-model ref="editForm" :model="model" :rules="rules" :label-col="labelCol" :wrapper-col="wrapperCol">
      <a-form-model-item :label="$t('form.name')" prop="field" :wrapper-col="wrapperColHalf">
        <a-input v-model="model.field" />
      </a-form-model-item>

      <a-form-model-item :label="$t('form.prefix')" :wrapper-col="wrapperColHalf">
        <div class="input-group">
          <a-form-model-item prop="prefix">
            <a-input v-model="model.prefix" />
          </a-form-model-item>
          <span class="input-group-addon">{{$t('form.postfix')}}</span>
          <a-form-model-item prop="postfix">
            <a-input v-model="model.postfix" />
          </a-form-model-item>
          <span class="input-group-addon">{{$t('form.divider')}}</span>
          <a-form-model-item prop="divider">
            <a-input v-model="model.divider" />
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.loop')" :wrapper-col="wrapperColHalf">
        <div class="input-group">
          <a-form-model-item prop="loop">
            <a-input v-model="model.loop" :placeholder="$t('tips.range.int')" />
          </a-form-model-item>
          <span class="input-group-addon">{{$t('form.loopfix')}}</span>
          <a-form-model-item prop="loopfix">
            <a-input v-model="model.loopfix" />
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.type')" :wrapper-col="wrapperColHalf">
        <div class="input-group">
          <a-form-model-item prop="type">
            <a-select v-model="model.type">
              <a-select-option value="list">{{$t('form.type.list')}}</a-select-option>
              <a-select-option value="timestamp">{{$t('form.type.timestamp')}}</a-select-option>
            </a-select>
          </a-form-model-item>
          <span class="input-group-addon"> {{$t('form.mode')}}</span>
          <a-form-model-item prop="mode">
            <a-select v-model="model.mode">
              <a-select-option value="parallel">{{$t('form.mode.parallel')}}</a-select-option>
              <a-select-option value="recursive">{{$t('form.mode.recursive')}}</a-select-option>
            </a-select>
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.width')" :wrapper-col="wrapperColHalf">
        <div class="input-group">
          <a-form-model-item prop="length">
            <a-input v-model="model.length" :min="0" />
          </a-form-model-item>
          <span class="input-group-addon">{{$t('form.left.pad')}}</span>
          <a-form-model-item prop="leftPad">
            <a-input v-model="model.leftPad" />
          </a-form-model-item>
          <span class="input-group-addon">{{$t('form.right.pad')}}</span>
          <a-form-model-item prop="rightPad">
            <a-input v-model="model.rightPad" />
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.format')" prop="format" :wrapperCol="wrapperColHalf">
        <div class="input-group">
          <a-input v-model="model.format">
            <a-select v-model="model.format" slot="addonAfter" default-value="" style="width: 100px">
              <a-select-option value="">
                {{$t('form.function')}}
              </a-select-option>
              <a-select-option value="md5">md5</a-select-option>
              <a-select-option value="sha1">sha1</a-select-option>
              <a-select-option value="base64">base64</a-select-option>
              <a-select-option value="urlencode">urlencode</a-select-option>
            </a-select>
          </a-input>
          <a-form-model-item :label="$t('form.rand')" prop="rand" :wrapperCol="wrapperColHalf">
            <a-switch v-model="model.rand" />
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.desc')" prop="note">
        <a-input v-model="model.note" type="textarea" rows="3" />
      </a-form-model-item>

      <a-form-model-item class="center" :wrapper-col="{ span: 19, offset: 4 }">
        <a-button @click="save" type="primary">{{$t('form.save')}}</a-button>
        <a-button @click="reset" style="margin-left: 10px;">{{$t('form.reset')}}</a-button>
      </a-form-model-item>
    </a-form-model>
  </div>
</template>

<script>
import {saveDefField, saveInstancesItem} from "../api/manage";
import {checkLoop} from "../api/utils";
import {labelColFull, wrapperColFull} from "../utils/const";

export default {
  name: 'FieldInfoComponent',
  data() {
    return {
      labelCol: labelColFull,
      wrapperCol: wrapperColFull,

      colsFull: 24,
      colsHalf: 12,
      labelColFull: { lg: { span: 4 }, sm: { span: 4 } },
      wrapperColFull: { lg: { span: 16 }, sm: { span: 16 } },
      labelColHalf: { lg: { span: 8}, sm: { span: 8 } },
      labelColHalf2: { lg: { span: 4}, sm: { span: 4 } },
      wrapperColHalf: { lg: { span: 10 }, sm: { span: 10 } },
      rules: {
        field: [
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
        // loop: [
        //   { validator: checkLoop, message: this.$i18n.t('valid.loop.check'), trigger: 'change' },
        // ],
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
