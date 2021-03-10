<template>
  <div>
    <a-form-model ref="editForm" :model="model" :rules="rules">
      <a-form-model-item :label="$t('form.name')" prop="title" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <a-input v-model="model.title" />
      </a-form-model-item>

      <a-form-model-item :label="$t('form.dir')" prop="folder" class="zui-input-group zui-input-with-tips"
                          :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <div class="input-group">
          <a-form-model-item prop="folder">
            <a-select v-model="model.folder">
              <a-select-option v-for="(item, index) in dirs" :value="item.name" :key="index">
                {{item.name}}</a-select-option>
            </a-select>
          </a-form-model-item>
          <span class="input-group-addon">{{ $t('form.folder') }}</span>
          <a-form-model-item :style="{ display: 'inline-block', width: 'calc(30% - 30px)' }">
            <a-input v-model="model.subFolder"></a-input>
          </a-form-model-item>
        </div>
      </a-form-model-item>

      <a-form-model-item :label="$t('form.file.name')" prop="fileName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
        <a-input v-model="model.fileName" />
      </a-form-model-item>

    <a-form-model-item :label="$t('form.prefix')" prop="prefix" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
      <div class="input-group">
        <a-input v-model="model.prefix" />
        <span class="input-group-addon">{{ $t('form.postfix') }}</span>
        <a-form-model-item prop="postfix">
          <a-input v-model="model.postfix" />
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
import {getInstances, saveInstances} from "../../../../api/manage";
import {checkDirIsYaml} from "../../../../api/utils";
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
  name: 'RangesEdit',
  props: {
    afterSave: Function,
    id: {
      type: [Number, String],
      default: function() {
        return this.$route.params.id;
      }
    },
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
        folder: [
          { validator: checkDirIsYaml, trigger: 'change' },
        ],
      },

      model: {folder: 'yaml/'},
      dirs: [],
      workDir: '',
    };
  },
  watch: {
    id: function(newId, oldId) {
      if (newId == oldId) {
        return;
      }
      this.loadData();
    }
  },
  mounted () {
    this.loadData();
  },
  methods: {
    loadData () {
      let id = this.id;
      if (id === null) {
        return;
      }
      if (id) {
        if (typeof id === 'string') id = Number.parseInt(id);
        getInstances(id).then(json => {
          console.log('getInstances', json)
          this.model = json.data
          this.dirs = json.res
          this.workDir = json.workDir
        })
      } else {
        this.reset();
      }
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
        saveInstances(this.model).then(json => {
          console.log('saveInstances', json)
          if (this.afterSave) {
            this.afterSave(json);
          }
        })
      })
    },
    reset() {
      console.log('reset')
      this.model = {folder: 'yaml/'};
      this.$refs.editForm.reset()
    },
  }
}
</script>

<style lang="less" scoped>
</style>
