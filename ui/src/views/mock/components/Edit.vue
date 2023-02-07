<template>
  <div class="mock-edit-modal">
    <a-modal
      :title="$t('msg.mock.create')"
      :visible="visible"
      :closable=true
      :footer="null"
      @cancel="cancel"
      width="100%"
      dialogClass="full-screen-modal"
    >
      <div class="mock-edit-main">
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
            </div>
            <div class="upload-content">
              <pre>{{ specContent }}</pre>
            </div>
          </a-col>

          <a-col :span="13" class="content-col">
            <a-tabs default-active-key="1" :animated="false">
              <a-tab-pane key="1" :tab="$t('msg.mock.mock')">
                <pre>{{ mockContent }}</pre>
              </a-tab-pane>
              <a-tab-pane key="2" :tab="$t('msg.mock.data')">
                <pre>{{ dataContent }}</pre>
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

export default {
  name: 'MockEditComp',
  components: {
  },
  data() {
    const styl = 'height: ' + (document.documentElement.clientHeight - 56) + 'px;'
    return {
      styl: styl,
      modelData: {},

      specContent: null,
      mockContent: null,
      dataContent: null,
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
    model: {
      type: Object,
      default: () => {return {}}
    },
    time: {
      type: Number,
      default: () => 0
    },
  },

  computed: {
  },
  created () {
    console.log('created')
    this.loadData()

    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadData()
    })
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
      this.loadData()
    },
    cancel() {
      console.log('cancel')
      this.$emit('cancel')
    },

    loadData () {
      console.log('loadData', this.modelProp)

      if (!this.model?.id) return

      // getDefFieldTree(this.modelProp.id).then(json => {
      //   console.log('getDefFieldTree', json)
      //   this.loadTreeCallback(json, selectedKey)
      // })
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
          this.specContent = json.data.spec
          this.mockContent = json.data.mock
          this.dataContent = json.data.data
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
