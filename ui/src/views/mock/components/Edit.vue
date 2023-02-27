<template>
  <div class="mock-edit-modal">
    <a-modal
        :closable=false
        :footer="null"
        :title="model.id == undefined ? $t('msg.mock.create') : $t('msg.mock.edit')"
        :visible="visible"
        dialogClass="full-screen-modal"
        width="100%"
    >
      <div class="mock-edit-main">
        <div class="buttons">
          <a-button :disabled="!readyToSave" type="primary" @click="save">
            {{ $t('form.save') }}
          </a-button> &nbsp;&nbsp;&nbsp;
          <a-button @click="cancel">
            {{ $t('form.close') }}
          </a-button>
        </div>

        <a-row :gutter="10" class="content-row">
          <a-col v-if="model.id == undefined" :span="11" class="content-col">
            <div class="upload-bar">
              <a-row>
                <a-col :span="5">
                  <a-upload :before-upload="beforeUpload"
                            :showUploadList="false"
                            accept=".yaml,.yml,.json">
                    <a-button>
                      <a-icon type="upload"/>
                      <span>{{ $t('upload.spec') }}</span>
                    </a-button>
                  </a-upload>
                </a-col>
                <a-col :span="9">
                  <span class="title">{{ model.name }}</span>
                </a-col>
                <a-col :span="10">
                  <span class="label-path"></span>
                  <a-input v-model="model.path" :placeholder="$t('msg.mock.input.path')"/>
                </a-col>
              </a-row>

            </div>
            <div class="upload-content">
              <pre>{{ model.specContent }}</pre>
            </div>
          </a-col>

          <a-col :span="model.id === undefined ? 13 : 24" class="content-col">
            <a-tabs :activeKey="currentTab" :animated="false" @change="tabChange">
              <a-tab-pane key="data" :tab="$t('msg.mock.data')">
                <pre v-show="model.id === undefined">{{ model.dataContent }}</pre>
                <design-in-component
                    v-show="model.id !== undefined"
                    ref="designPage"
                    :modelProp="dataModel"
                    :time="time"
                    :type="resType"
                    @save="onModelDataSave">
                    :visible="true">
                </design-in-component>
              </a-tab-pane>
              <a-tab-pane key="mock" :tab="$t('msg.mock.mock')">
                <!--                <pre>{{ model.mockContent }}</pre>-->
                <div class="yaml-editor">
                  <yaml-editor v-model="model.mockContent"/>
                </div>
              </a-tab-pane>
            </a-tabs>
          </a-col>
        </a-row>

      </div>
    </a-modal>

  </div>
</template>

<script>
import {uploadMock} from "@/api/mock";
import mockMixin from "@/store/mockMixin";
import {DesignInComponent} from '../../../components'
import {ResTypeDef} from "../../../api/utils";
import YamlEditor from './Yaml.vue';
import {getMock} from "@/api/mock";


export default {
  name: 'MockEditComp',
  components: {
    DesignInComponent,
    YamlEditor,
  },
  data() {
    return {
      model: {mockContent:''},
      resType: ResTypeDef,
      specReady: false,
      currentTab: this.current,
      cmOptions: {
        lineNumbers: true, // 显示行号
        mode: 'text/x-yaml', // 语法model
        gutters: ['CodeMirror-lint-markers'],  // 语法检查器
        theme: 'monokai', // 编辑器主题
        lint: true // 开启语法检查
      },
    };
  },
  props: {
    type: {
      type: String,
      required: true
    },
    visible: {
      type: Boolean,
      required: true
    },
    mock: {
      type: Object,
      default: () => null
    },
    time: {
      type: Number,
      default: () => 0
    },
    current: {
      type: String,
      default: () => ''
    },
  },
  mixins: [mockMixin],
  computed: {
    readyToSave() {
      return this.specReady && this.model.path?.trim()
    },
    codemirror() {
      return this.$refs.cmEditor.codemirror
    },
    dataModel() {
        return {id: this.model.defId}
    },
  },
  created() {
    console.log('created')
  },
  mounted: function () {
    console.log('mounted')
  },
  beforeDestroy() {
    console.log('beforeDestroy')
  },
  watch: {
    mock(val) {
      if(val == undefined){
        val = {mockContent:''};
      }
      this.model = val
      if(val.id !== undefined){
        this.specReady = true;
      }
      console.log("watch mock :", val)
    },
    current(val) {
      this.currentTab = val
    }
  },

  methods: {
    save() {
      console.log('save')
      this.saveMockItem(this.model).then((json) => {
        console.log('saveMockItem', json)
        if (json.code === 0) {
          this.model = {mockContent:''}
          this.specReady = false
          this.$emit('ok')
        } else {
          this.$notification['warning']({
            message: this.$i18n.t('upload.spec.failed'),
            duration: 3,
          });
        }
      })
    },
    cancel() {
      console.log('cancel')
      this.model = {mockContent:''}
      this.specReady = false
      this.$emit('cancel')
    },
    tabChange(key) {
      this.currentTab = key
    },

    getModel(id) {
      console.log('getModel', id)
    },

    onModelDataSave() {
      getMock(this.model.id).then(json =>{
          this.model.dataContent = json.data.dataContent;
      })
    },

    beforeUpload(file) {
      console.log('beforeUpload', file)

      const formData = new FormData()
      formData.append('file', file)

      uploadMock(formData).then((json) => {
        console.log('uploadMock', json)
        if (json.code === 0) {
          this.model = {
            name: json.data.name,
            specContent: json.data.spec,
            mockContent: json.data.mock,
            dataContent: json.data.data,
            dataPath: json.data.dataPath,
          }

          this.specReady = true

        } else {
          this.$notification['warning']({
            message: this.$i18n.t('upload.spec.failed'),
            duration: 3,
          });
        }
      })

      return false
    },
    handleEditSave() {
    },
  }
}
</script>

<style lang="less" scoped>
.mock-edit-main {

  .buttons {
    position: absolute;
    top: 6px;
    right: 6px;
    padding: 5px;
  }

}

</style>

<style lang="less">
.ant-modal-content {
  overflow: hidden;

  .ant-modal-body {
    height: calc(~ "100% - 55px");

    .mock-edit-main {
      height: 100%;

      .content-row {
        height: 100%;

        .content-col {
          height: 100%;
          .yaml-editor {
            height: 100%;
          }

          pre {
            padding: 10px;
            height: 100%;
          }

          .upload-bar {
            padding: 10px;

            .title {
              display: inline-block;
              padding: 3px 16px;
              font-size: larger;
              font-weight: bolder;
            }

            .label-path:before {
              display: inline-block;
              margin-right: 4px;
              color: #f5222d;
              font-size: 14px;
              font-family: SimSun, sans-serif;
              line-height: 1;
              content: "*";
            }

            input {
              width: calc(~ "100% - 20px");
            }

          }

          .upload-content {
            height: calc(~ "100% - 50px");
          }

          .ant-tabs {
            height: 100%;

            .ant-tabs-bar {
              margin-bottom: 10px;
            }

            .ant-tabs-content {
              height: calc(~ "100% - 40px");
              overflow-y: auto;

              .ant-tabs-tabpane-active {
                height: calc(~ "100% - 10px");
              }

            }
          }
        }
      }
    }
  }
}
</style>

<style lang="less">
.CodeMirror pre.CodeMirror-line {
  padding: 3px !important;
}
</style>
