<template>
  <div>
    <div class="head">
      <div class="title">
        {{$t('menu.instances.item.edit')}}
      </div>
      <div class="buttons"></div>
    </div>
    <div>
      <a-form-model ref="editForm" :model="model" :rules="rules">
        <a-row :gutter="colsFull">
            <a-form-model-item :label="$t('form.name')" prop="name" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
              <a-input v-model="model.name" />
            </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.value')" prop="value" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.value" />
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
  </div>
</template>

<script>
import {saveInstancesItem} from "../api/manage";
import {colsFull, colsHalf, labelColFull, wrapperColFull} from "@/utils/const";

export default {
  name: 'ResInstancesItemComponent',
  data() {
    return {
      colsFull: colsFull,
      colsHalf: colsHalf,
      labelColFull: labelColFull,
      wrapperColFull: wrapperColFull,

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

        saveInstancesItem(this.model).then(json => {
          console.log('saveInstancesItem', json)
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
