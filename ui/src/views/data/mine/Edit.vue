<template>
  <div>
    <div class="head">
      <div class="title">
        <span v-if="id==0">{{ $t('title.data.create') }}</span>
        <span v-if="id!=0">{{ $t('menu.data.edit') }}</span>
      </div>
      <div class="filter"></div>
      <div class="buttons">
        <a-button type="primary" @click="back()">{{ $t('action.back') }}</a-button>
      </div>
    </div>

    <div>
      <a-form-model ref="editForm" :model="model" :rules="rules" :label-col="labelCol" :wrapper-col="wrapperCol">
        <a-form-model-item :label="$t('form.name')" prop="title">
          <a-input v-model="model.title" />
        </a-form-model-item>

        <a-form-model-item :label="$t('form.dir')" class="zui-input-group zui-input-with-tips">
          <a-form-model-item prop="folder" :style="{ display: 'inline-block', width: 'calc(70% - 30px)' }">
            <a-select v-model="model.folder">
              <a-select-option v-for="(item, index) in dirs" :value="item.name" :key="index">
                {{item.name}}</a-select-option>
            </a-select>
          </a-form-model-item>

          <span class="zui-input-group-addon" :style="{ width: '60px' }">
            <span>{{ $t('form.folder') }}</span>
          </span>

          <a-form-model-item :style="{ display: 'inline-block', width: 'calc(30% - 30px)' }">
            <a-input v-model="model.subFolder"></a-input>
          </a-form-model-item>
        </a-form-model-item>

        <a-form-model-item :label="$t('form.file.name')" prop="fileName">
          <a-input v-model="model.fileName" />
        </a-form-model-item>

        <a-form-model-item :label="$t('form.desc')" prop="desc">
          <a-input v-model="model.desc" type="textarea" rows="3" />
        </a-form-model-item>

        <a-form-model-item :wrapper-col="{ span: 14, offset: 6 }">
          <a-button @click="save" type="primary">
            {{$t('action.save')}}
          </a-button>
          <a-button @click="reset" style="margin-left: 10px;">
            {{$t('action.reset')}}
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
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
        fileName: [
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
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
