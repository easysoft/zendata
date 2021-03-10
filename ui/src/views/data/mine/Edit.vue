<template>
  <div>
    <a-form-model ref="editForm" :model="model" :rules="rules" :label-col="labelCol" :wrapper-col="wrapperCol" :colon="false">
      <a-form-model-item :label="$t('form.name')" prop="title">
        <a-input v-model="model.title" />
      </a-form-model-item>

      <a-form-model-item :label="$t('form.dir')" prop="folder">
        <div class="input-group">
          <a-select v-model="model.folder">
            <a-select-option v-for="(item, index) in dirs" :value="item.name" :key="index">
              {{item.name}}</a-select-option>
          </a-select>
          <span class="input-group-addon fix-border">{{ $t('form.folder') }}</span>
          <a-input v-model="model.subFolder"></a-input>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.file.name')" prop="fileName">
        <a-input v-model="model.fileName" />
      </a-form-model-item>

      <a-form-model-item :label="$t('form.desc')" prop="desc">
        <a-input v-model="model.desc" type="textarea" rows="3" />
      </a-form-model-item>

      <a-form-model-item :wrapper-col="{ span: 18, offset: 4 }">
        <a-button @click="save" type="primary">
          {{$t('form.save')}}
        </a-button>
        <a-button @click="reset">
          {{$t('form.reset')}}
        </a-button>
      </a-form-model-item>
    </a-form-model>
  </div>
</template>

<script>

import {getDef, saveDef} from "../../../api/manage";
import {labelCol, wrapperCol} from "../../../utils/const";
import {checkDirIsUsers} from "../../../api/utils";

export default {
  name: 'Mine',
  props: {
    id: {
      type: Number,
      default: function() {
        return this.$route.params.id;
      }
    },
    afterSave: Function
  },
  data() {
    return {
      labelCol: labelCol,
      wrapperCol: wrapperCol,
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
      model: { folder: 'users/', type: 'text' },
      dirs: [],
      workDir: '',
    };
  },
  watch: {
    id: function() {
      console.log('watch id ' + this.id)
      this.loadData();
    }
  },
  mounted () {
    this.loadData();
  },
  methods: {
    loadData() {
      if (this.id === null) {
        return;
      }
      if (this.id) {
        getDef(this.id).then(json => {
          this.model = json.data
          this.dirs = json.res
          this.workDir = json.workDir
        })
      } else {
        this.reset();
      }
    },
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
          if (this.afterSave) {
            this.afterSave(json);
          }
        })
      })
    },
    reset () {
      this.$refs.editForm.resetFields()
    },
  }
}
</script>

<style lang="less" scoped>

</style>
