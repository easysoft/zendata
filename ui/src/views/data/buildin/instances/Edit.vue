<template>
  <div>
    <div class="head">
      <div class="title">
        序列<span v-if="id!=0">编辑</span><span v-if="id==0">新建</span>
      </div>
      <div class="buttons">
        <a-button type="primary" @click="back()">返回</a-button>
      </div>
    </div>

    <div>
      <a-form-model ref="editForm" :model="model" :rules="rules">
        <a-row :gutter="colsFull">
          <a-form-model-item label="名称" prop="title" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.title" />
          </a-form-model-item>
        </a-row>
        <a-row :gutter="colsFull">
          <a-form-model-item label="引用" prop="name" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            {{model.name}}
          </a-form-model-item>
        </a-row>

      <a-row :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="前缀" prop="prefix" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="model.prefix" />
          </a-form-model-item>
        </a-col>
        <a-col :span="colsHalf">
          <a-form-model-item label="后缀" prop="postfix" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
            <a-input v-model="model.postfix" />
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item label="格式" prop="format" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <div class="inline">
              <a-input v-model="model.format">
                <a-select slot="addonAfter" default-value="" style="width: 80px">
                  <a-select-option value="">
                    函数
                  </a-select-option>
                  <a-select-option value=".jp">
                    md5
                  </a-select-option>
                </a-select>
              </a-input>
            </div>
          </a-form-model-item>
        </a-col>
        <a-col :span="colsHalf"></a-col>
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
  </div>
</template>

<script>
import {getInstances, saveInstances} from "../../../../api/manage";

export default {
  name: 'RangesEdit',
  data() {
    return {
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
      },

      id: 0,
      model: {},
    };
  },

  computed: {
  },
  created () {
    this.id = parseInt(this.$route.params.id)
    console.log(this.id)
    this.loadData()
  },
  mounted () {

  },
  methods: {
    loadData () {
      if (!this.id) return

      getInstances(this.id).then(res => {
        console.log('getInstances', res)
        this.model = res.data
      })
    },
    save() {
      console.log('save')
      this.$refs.editForm.validate(valid => {
        console.log(valid, this.model)
        if (!valid) {
          console.log('validation fail')
          return
        }

        saveInstances(this.model).then(json => {
          console.log('saveInstances', json)
          this.back()
        })
      })
    },
    reset() {
      console.log('reset')
      this.$refs.editForm.reset()
    },
    back() {
      this.$router.push({path: '/data/buildin/instances/list'});
    },
  }
}
</script>

<style lang="less" scoped>
</style>
