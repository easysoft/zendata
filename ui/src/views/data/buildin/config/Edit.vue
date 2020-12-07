<template>
  <div>
    <div class="head">
      <div class="title">
        <div class="title">
          <span v-if="id==0">{{ $t('title.config.create') }}</span>
          <span v-if="id!=0">{{ $t('menu.config.edit') }}</span>
        </div>
      </div>
      <div class="filter"></div>
      <div class="buttons">
        <a-button type="primary" @click="back()">{{ $t('action.back') }}</a-button>
      </div>
    </div>

    <div>
      <a-form-model ref="editForm" :model="model" :rules="rules">
        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.name')" prop="title" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.title" />
          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.dir')" class="zui-input-group zui-input-with-tips"
                             :labelCol="labelColFull" :wrapperCol="wrapperColFull">

            <a-form-model-item prop="folder" :style="{ display: 'inline-block', width: 'calc(70% - 40px)' }">
              <a-select v-model="model.folder">
                <a-select-option v-for="(item, index) in dirs" :value="item.name" :key="index">
                  {{item.name}}</a-select-option>
              </a-select>
            </a-form-model-item>

            <span class="zui-input-group-addon" :style="{ width: '80px' }">
              <span>{{ $t('form.folder') }}</span>
            </span>

            <a-form-model-item :style="{ display: 'inline-block', width: 'calc(30% - 40px)' }">
              <a-input v-model="model.subFolder"></a-input>
            </a-form-model-item>

          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.file.name')" prop="fileName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.fileName" />
          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.prefix')" class="zui-input-group zui-input-with-tips"
                             :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-form-model-item prop="prefix" :style="{ display: 'inline-block', width: 'calc(70% - 40px)' }">
              <a-input v-model="model.prefix" />
            </a-form-model-item>

            <span class="zui-input-group-addon" :style="{ width: '80px' }">
              <span>{{$t('form.postfix')}}</span>
            </span>

            <a-form-model-item prop="postfix" :style="{ display: 'inline-block', width: 'calc(30% - 40px)' }">
              <a-input v-model="model.postfix" />
            </a-form-model-item>
          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.loop')" class="zui-input-group zui-input-with-tips"
                             :labelCol="labelColFull" :wrapperCol="wrapperColFull">

            <a-form-model-item prop="loop" :style="{ display: 'inline-block', width: 'calc(70% - 40px)' }">
              <a-input v-model="model.loop" placeholder="数字或数字区间" />
            </a-form-model-item>

            <span class="zui-input-group-addon" :style="{ width: '80px' }">
              <span>{{$t('form.loopfix')}}</span>
            </span>

            <a-form-model-item prop="loopfix" :style="{ display: 'inline-block', width: 'calc(30% - 40px)' }">
              <a-input v-model="model.loopfix" />
            </a-form-model-item>
          </a-form-model-item>
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.format')" class="zui-input-group zui-input-with-tips"
                             :labelCol="labelColFull" :wrapperCol="wrapperColFull">

            <a-form-model-item prop="format" :style="{ display: 'inline-block', width: 'calc(70% - 40px)' }">
              <a-input v-model="model.format"></a-input>
            </a-form-model-item>

            <span class="zui-input-group-addon" :style="{ width: '80px' }">
              <span>{{$t('form.function')}}</span>
            </span>

            <a-form-model-item prop="format" :style="{ display: 'inline-block', width: 'calc(30% - 40px)' }">
              <a-select v-model="model.format">
                <a-select-option value="md5">md5</a-select-option>
                <a-select-option value="sha1">sha1</a-select-option>
                <a-select-option value="base64">base64</a-select-option>
                <a-select-option value="urlencode">urlencode</a-select-option>
              </a-select>
            </a-form-model-item>

          </a-form-model-item>
        </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item :label="$t('form.desc')" prop="desc" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-input v-model="model.desc" type="textarea" rows="3" />
        </a-form-model-item>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item class="center">
          <a-button @click="save" type="primary">{{$t('action.save')}}</a-button>
          <a-button @click="reset" style="margin-left: 10px;">{{$t('action.reset')}}</a-button>
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
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
        fileName: [
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
        loop: [
          { validator: checkLoop, message: this.$i18n.t('valid.loop.format'), trigger: 'change' },
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
