<template>
  <div>
    <div class="head">
      <div class="title">
        <div class="title">
          <span v-if="id==0">{{ $t('title.instances.create') }}</span>
          <span v-if="id!=0">{{ $t('menu.instances.edit') }}</span>
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
          <a-form-model-item :label="$t('form.dir')" prop="folder" class="zui-input-group zui-input-with-tips"
                             :labelCol="labelColFull" :wrapperCol="wrapperColFull">
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
        </a-row>

        <a-row :gutter="colsFull">
          <a-form-model-item :label="$t('form.file.name')" prop="fileName" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
            <a-input v-model="model.fileName" />
          </a-form-model-item>
        </a-row>

      <a-row :gutter="colsFull">
        <a-col :span="colsHalf">
          <a-form-model-item :label="$t('form.prefix')" prop="prefix" :labelCol="labelColHalf" :wrapperCol="wrapperColHalf">
            <a-input v-model="model.prefix" />
          </a-form-model-item>
        </a-col>
        <a-col :span="colsHalf">
          <a-form-model-item :label="$t('form.postfix')" prop="postfix" :labelCol="labelColHalf2" :wrapperCol="wrapperColHalf">
            <a-input v-model="model.postfix" />
          </a-form-model-item>
        </a-col>
      </a-row>

      <a-row :gutter="colsFull">
        <a-form-model-item :label="$t('form.desc')" prop="desc" :labelCol="labelColFull" :wrapperCol="wrapperColFull">
          <a-input v-model="model.desc" type="textarea" rows="3" />
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

      id: 0,
      model: {folder: 'yaml/'},
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
      if (!this.id) return

      getInstances(this.id).then(json => {
        console.log('getInstances', json)
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
