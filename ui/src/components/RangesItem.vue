<template>
  <div>
    <div class="head">
      <div class="title">
        {{$t('menu.ranges.item.edit')}}
      </div>
      <div class="buttons"></div>
    </div>
    <div>
      <a-form-model ref="editForm" :model="model" :rules="rules">
        <a-row :gutter="colsFull">
            <a-form-model-item :label="$t('form.name')" prop="field" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
              <a-input v-model="model.field" />
            </a-form-model-item>
        </a-row>
        <a-row :gutter="colsFull">
          <a-form-model-item class="center">
            <a-button @click="save" type="primary">{{$t('form.save')}}</a-button>
            <a-button @click="reset" style="margin-left: 10px;">{{$t('form.reset')}}</a-button>
          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-col :offset="3">
            <field-range-component
                ref="rangeComp"
                :type="'ranges'"
                :model="model"
                :time2="time">
            </field-range-component>
          </a-col>
          <a-col :span="2"></a-col>
        </a-row>
      </a-form-model>
    </div>
  </div>
</template>

<script>
import {saveRangesItem} from "../api/manage";
import FieldRangeComponent from "./FieldRange";

export default {
  name: 'ResRangesItemComponent',
  components: {FieldRangeComponent},
  data() {
    return {
      colsFull: 24,
      colsHalf: 12,
      labelColFull: { lg: { span: 4 }, sm: { span: 4 } },
      wrapperColFull: { lg: { span: 16 }, sm: { span: 16 } },
      rules: {
        field: [
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
      },
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

        saveRangesItem(this.model).then(json => {
          console.log('saveRangesItem', json)
          this.$emit('save')
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
</style>
