<template>
  <div>
    <div class="head">
      <div class="title">
        测试数据<span v-if="id!=0">编辑</span><span v-if="id==0">新建</span>
      </div>
      <div class="filter"></div>
      <div class="buttons">
        <a-button type="primary" @click="back()">返回</a-button>
      </div>
    </div>

    <div>
      <a-form-model ref="editForm" :model="model" :rules="rules" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-model-item label="名称" prop="title">
          <a-input v-model="model.title" />
        </a-form-model-item>

        <a-form-model-item label="目录" class="zui-input-group zui-input-with-tips">
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

        <a-form-model-item label="文件名" prop="fileName">
          <a-input v-model="model.fileName" />
        </a-form-model-item>

<!--        <a-form-model-item label="类型" prop="type">
          <a-select v-model="model.type">
            <a-select-option value="text">字符串</a-select-option>
            <a-select-option value="article">文章</a-select-option>
          </a-select>
        </a-form-model-item>-->

        <a-form-model-item label="描述" prop="desc">
          <a-input v-model="model.desc" type="textarea" rows="3" />
        </a-form-model-item>

        <a-form-model-item :wrapper-col="{ span: 14, offset: 6 }">
          <a-button @click="save" type="primary">
            保存
          </a-button>
          <a-button @click="reset" style="margin-left: 10px;">
            重置
          </a-button>
        </a-form-model-item>
      </a-form-model>
    </div>
  </div>
</template>

<script>

import { getDef, saveDef } from "../../../api/manage";
import { labelColLarge, wrapperColLarge } from "../../../utils/const";
import {checkDirIsUsers} from "../../../api/utils";

export default {
  name: 'Mine',
  data() {
    return {
      labelCol: labelColLarge,
      wrapperCol: wrapperColLarge,
      rules: {
        title: [
          { required: true, message: '名称不能为空', trigger: 'change' },
        ],
        fileName: [
          { required: true, message: '文件名不能为空', trigger: 'change' },
        ],
        folder: [
          { validator: checkDirIsUsers, trigger: 'change' },
        ],
      },
      id: 0,
      model: { folder: 'users/', type: 'text' },
      dirs: [],
      workDir: '',
    };
  },
  computed: {

  },
  created () {
    this.id = parseInt(this.$route.params.id)
    console.log(this.id)
    if (this.id == 0) return

    getDef(this.id).then(json => {
      console.log('getDef', json)
      this.model = json.data
      this.dirs = json.res
      this.workDir = json.workDir
    })
  },
  mounted () {

  },
  methods: {
    save() {
      this.$refs.editForm.validate(valid => {
        console.log(valid, this.model)
        if (!valid) {
          console.log('validation fail')
          return
        }

        if (this.model.subFolder && this.model.subFolder != '') this.model.folder += this.model.subFolder
        saveDef(this.model).then(json => {
          console.log('saveDef', json)
          this.back()
        })
      })
    },
    reset () {
      this.$refs.editForm.resetFields()
    },
    back() {
      this.$router.push({path: '/data/mine/list'});
    },
  }
}
</script>

<style scoped>

</style>
