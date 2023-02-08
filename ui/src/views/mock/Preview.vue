<template>
  <div class="mock-preview-main">
    <a-card>
      <span slot="title">
        <a-icon type="profile" />
      </span>

      <div>
        <div v-for="(path, url) in mockItem.paths" :key="url" class="path">
          <div>{{url}}</div>
          <div v-for="(methodVal, method) in path" :key="method" class="method">
            <div>{{method}}</div>
            <div v-for="(codeVal, code) in methodVal" :key="code" class="code">
              <div>{{code}}</div>
              <div v-for="(mediaVal, media) in codeVal" :key="media" class="media">
                <a-popover
                    :title="$t('msg.mock.response')"
                    trigger="click"
                    :visible="clicked"
                    @visibleChange="handleClickChange"
                >
                  <div slot="content">
                    <div class="mock-preview-resp">
                      <pre>{{respSample}}</pre>
                    </div>
                  </div>
                  <a @click="preview(url, method, code, media)">{{media}}</a>
                </a-popover>
              </div>
            </div>
          </div>
        </div>

      </div>

    </a-card>
  </div>
</template>

<script>

import mockMixin from "@/store/mockMixin";
import {getPreviewData} from "@/api/mock";

export default {
  name: 'MockPreview',
  components: {
  },
  props: {
  },
  mixins: [mockMixin],
  data: function() {
    return {
      clicked: false,
      hovered: false,
      respSample: null,
    };
  },
  mounted: function() {
  },
  methods: {
    preview(url, method, code, media) {
      console.log(url, method, code, media)

      getPreviewData(this.mockItem.id, url, method, code, media).then(json => {
        if (json.code === 0) {
          this.respSample = json.data
        }
      })
    },

    handleClickChange(visible) {
      this.clicked = visible;
      this.hovered = false;
    },
    hide() {
      this.clicked = false;
      this.hovered = false;
    },
  },
  watch: {

  }
}
</script>

<style lang="less" scoped>
.mock-preview-main {
  .path {
    padding-left: 10px;
    .method {
      padding-left: 10px;
      .code {
        padding-left: 10px;
        .media {
          padding-left: 10px;
        }
      }
    }
  }
}
</style>

<style lang="less">
.mock-preview-resp {
  max-width: 600px;
  word-wrap:break-word;
}
</style>
