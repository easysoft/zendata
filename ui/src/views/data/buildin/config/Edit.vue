<template>
  <div>
    <div class="head">
      <div class="title">
        字段<span v-if="id!=0">编辑</span><span v-if="id==0">新建</span>
      </div>
      <div class="filter"></div>
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
          <a-form-model-item label="目录" prop="folder" class="zui-input-group zui-input-with-tips"
                             :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-form-model-item prop="folder" :style="{ display: 'inline-block', width: 'calc(70% - 30px)' }">
              <a-select v-model="model.folder" placeholder="请选择">
                <a-select-option v-for="(item, index) in dirs" :value="item.name" :key="index">
                  {{item.name}}</a-select-option>
              </a-select>
              <span class="zui-input-tips">工作目录：{{workDir}}</span>
            </a-form-model-item>

            <span class="zui-input-group-addon" :style="{ width: '60px' }">
              <span>子目录</span>
            </span>

            <a-form-model-item :style="{ display: 'inline-block', width: 'calc(30% - 30px)' }">
              <a-input v-model="model.subFolder"></a-input>
            </a-form-model-item>
          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item label="文件名" prop="fileName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.fileName" />
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
            <a-form-model-item label="循环" prop="loop" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
              <a-input v-model="model.loop" placeholder="数字3，或区间1-3" />
            </a-form-model-item>
          </a-col>
          <a-col :span="colsHalf">
            <a-form-model-item label="循环间隔符" prop="loopfix" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
              <a-input v-model="model.loopfix" />
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
        <a-form-model-item label="描述" prop="desc" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-input v-model="model.desc" type="textarea" rows="3" />
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
import {getConfig, saveConfig} from "../../../../api/manage";
import {checkLoop, checkDirIsYaml} from "../../../../api/utils";

export default {
  name: 'ConfigEdit',
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
        title: [
          { required: true, message: '名称不能为空', trigger: 'change' },
        ],
        fileName: [
          { required: true, message: '文件名不能为空', trigger: 'change' },
        ],
        loop: [
          { validator: checkLoop, message: '需为正整数或其区间', trigger: 'change' },
        ],
        folder: [
          { validator: checkDirIsYaml, trigger: 'change' },
        ],
      },

      id: 0,
      model: { folder: 'yaml/'},
      dirs: [],
      workDir: '',
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
      getConfig(this.id).then(json => {
        console.log('getConfig', json)
        this.model = json.data
        this.dirs = json.res
        this.workDir = json.workDir
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

        if (this.model.subFolder && this.model.subFolder != '') this.model.folder += this.model.subFolder
        saveConfig(this.model).then(json => {
          console.log('saveConfig', json)
          this.back()
        })
      })
    },
    reset() {
      console.log('reset')
      this.$refs.editForm.reset()
    },
    back() {
      this.$router.push({path: '/data/buildin/config/list'});
    },
  }
}
</script>

<style lang="less" scoped>
</style>
