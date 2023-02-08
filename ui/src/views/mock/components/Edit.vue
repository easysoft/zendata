<template>
  <div class="mock-edit-modal">
    <a-modal
      :title="$t('msg.mock.create')"
      :visible="visible"
      :closable=false
      :footer="null"
      width="100%"
      dialogClass="full-screen-modal"
    >
      <div class="mock-edit-main">
        <div class="buttons">
          <a-button @click="save" type="primary" :disabled="!readyToSave">
            {{ $t('form.save') }}
          </a-button> &nbsp;&nbsp;&nbsp;
          <a-button @click="cancel">
            {{ $t('form.close') }}
          </a-button>
        </div>

        <a-row :gutter="10" class="content-row">
          <a-col :span="11" class="content-col">
            <div class="upload-bar">
              <a-upload :before-upload="beforeUpload"
                        :showUploadList="false"
                        accept=".yaml,.yml,.json">
                <a-button>
                  <a-icon type="upload" />
                  <span>{{$t('upload.spec')}}</span>
                </a-button>
              </a-upload>

              <span class="title">{{model.name}}</span>
            </div>
            <div class="upload-content">
              <pre>{{ model.specContent }}</pre>
            </div>
          </a-col>

          <a-col :span="13" class="content-col">
            <a-tabs default-active-key="1" :animated="false">
              <a-tab-pane key="1" :tab="$t('msg.mock.mock')">
                <pre>{{ model.mockContent }}</pre>
              </a-tab-pane>
              <a-tab-pane key="2" :tab="$t('msg.mock.data')">
                <pre>{{ model.dataContent }}</pre>
              </a-tab-pane>
            </a-tabs>
          </a-col>
        </a-row>

      </div>
    </a-modal>

  </div>
</template>

<script>
import {} from "../../../api/manage";
import {uploadMock} from "@/api/mock";
import mockMixin from "@/store/mockMixin";

export default {
  name: 'MockEditComp',
  components: {
  },
  data() {
    return {
      model: {},
      readyToSave: false,
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
    time: {
      type: Number,
      default: () => 0
    },
  },
  mixins: [mockMixin],
  computed: {
  },
  created () {
    console.log('created')
  },
  mounted: function () {
    console.log('mounted')
  },
  beforeDestroy() {
    console.log('beforeDestroy')
  },

  methods: {
    save() {
      console.log('save')
      this.saveMockItem(this.model).then((json) => {
        console.log('saveMockItem', json)
        if (json.code === 0) {
          this.model = {}
          this.readyToSave = false
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
      this.model = {}
      this.readyToSave = false
      this.$emit('cancel')
    },

    getModel(id) {
      console.log('getModel', id)
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
          }

          this.readyToSave = true

        } else {
          this.$notification['warning']({
            message: this.$i18n.t('upload.spec.failed'),
            duration: 3,
          });
        }
      })

      return false
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
    height: calc(~"100% - 55px");

    .mock-edit-main {
      height: 100%;
      .content-row {
        height: 100%;

        .content-col {
          height: 100%;
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
          }
          .upload-content {
            height: calc(~"100% - 50px");
          }

          .ant-tabs {
            height: 100%;
            .ant-tabs-bar {
              margin-bottom: 10px;
            }
            .ant-tabs-content {
              height: calc(~"100% - 40px");
              overflow-y: auto;

              .ant-tabs-tabpane-active {
                height: calc(~"100% - 10px");
              }
            }
          }
        }
      }
    }
  }
}
</style>
