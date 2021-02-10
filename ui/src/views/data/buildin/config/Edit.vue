<template>
  <div>
    <a-form-model ref="editForm" :model="model" :rules="rules">
      <a-form-model-item :label="$t('form.name')" prop="title" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <a-input v-model="model.title" />
      </a-form-model-item>

      <a-form-model-item :label="$t('form.dir')"
                          :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <div class="input-group">
          <a-form-model-item prop="folder">
            <a-select v-model="model.folder">
              <a-select-option v-for="(item, index) in dirs" :value="item.name" :key="index">
                {{item.name}}</a-select-option>
            </a-select>
          </a-form-model-item>
          <span class="input-group-addon">{{ $t('form.folder') }}</span>
          <a-form-model-item>
            <a-input v-model="model.subFolder"></a-input>
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.file.name')" prop="fileName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <a-input v-model="model.fileName" />
      </a-form-model-item>

      <a-form-model-item :label="$t('form.prefix')"
                          :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <div class="input-group">
          <a-form-model-item prop="prefix">
            <a-input v-model="model.prefix" />
          </a-form-model-item>
          <span class="input-group-addon">{{$t('form.postfix')}}</span>
          <a-form-model-item prop="postfix">
            <a-input v-model="model.postfix" />
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.loop')"
                          :labelCol="labelColFull" :wrapperCol="wrapperColFull">

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

      <a-form-model-item :label="$t('form.format')"
                          :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <div class="input-group">
          <a-form-model-item prop="format">
            <a-input v-model="model.format"></a-input>
          </a-form-model-item>

          <span class="input-group-addon">{{$t('form.function')}}</span>

          <a-form-model-item prop="format">
            <a-select v-model="model.format">
              <a-select-option value="md5">md5</a-select-option>
              <a-select-option value="sha1">sha1</a-select-option>
              <a-select-option value="base64">base64</a-select-option>
              <a-select-option value="urlencode">urlencode</a-select-option>
            </a-select>
          </a-form-model-item>
        </div>
      </a-form-model-item>


    <a-form-model-item :label="$t('form.desc')" prop="desc" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
      <a-input v-model="model.desc" type="textarea" rows="3" />
    </a-form-model-item>

    <a-form-model-item class="center" :wrapper-col="{ span: 18, offset: 4 }">
      <a-button @click="save" type="primary">{{$t('form.save')}}</a-button>
      <a-button @click="reset" style="margin-left: 10px;">{{$t('form.reset')}}</a-button>
    </a-form-model-item>
  </a-form-model>
  </div>
</template>

<script>
import {getConfig, saveConfig} from "../../../../api/manage";
import {checkLoop, checkDirIsYaml} from "../../../../api/utils";
import {
  colsFull,
  colsHalf,
  labelColFull,
  labelColHalf,
  labelColHalf2,
  wrapperColFull,
  wrapperColHalf
} from "@/utils/const";

export default {
  name: 'ConfigEdit',
  props: {
    afterSave: Function
  },
  data() {
    return {
      colsFull: colsFull,
      colsHalf: colsHalf,
      labelColFull: labelColFull,
      wrapperColFull: wrapperColFull,
      labelColHalf: labelColHalf,
      labelColHalf2: labelColHalf2,
      wrapperColHalf: wrapperColHalf,

      rules: {
        title: [
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
        fileName: [
          { required: true, message: this.$i18n.t('valid.required'), trigger: 'change' },
        ],
        loop: [
          { validator: checkLoop, message: this.$i18n.t('valid.loop.check'), trigger: 'change' },
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
          console.log('saveConfig', json, this.afterSave)
          if (this.afterSave) {
            this.afterSave(json);
          }
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
