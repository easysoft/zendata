<template>
  <div class="mock-preview-main">
    <a-card>
      <span slot="title">
        <a-icon type="profile" />
      </span>

      <div v-if="mockItem">
        <div v-for="(path, url) in mockItem.paths" :key="url" class="path">
          <div>{{url}}</div>
          <div v-for="(methodVal, method) in path" :key="method" class="method">
            <div>{{method}}</div>
            <div v-for="(codeVal, code) in methodVal" :key="code" class="code">
              <div>{{code}}</div>
              <div v-for="(mediaVal, media) in codeVal" :key="media" class="media">
                <a @click="preview(mockItem.id, url, method, code, media)">{{media}}</a>
              </div>
            </div>
          </div>
        </div>

      </div>

    </a-card>

    <a-drawer
        :title="$t('msg.mock.response')"
        placement="left"
        :closable="false"
        :visible="responseVisible"
        width="50%"
        @close="closePreview">
      <div>
        <div class="mock-preview-resp">
          <pre>{{respSample}}</pre>
        </div>
      </div>
    </a-drawer>

  </div>
</template>

<script>

import mockMixin from "@/store/mockMixin";
import {getPreviewResp} from "@/api/mock";

export default {
  name: 'MockPreview',
  components: {
  },
  props: {
  },
  mixins: [mockMixin],
  data: function() {
    return {
      responseVisible: false,
      hovered: false,
      respSample: null,
    };
  },
  mounted: function() {
  },
  methods: {
    preview(id, url, method, code, media) {
      console.log(id, url, method, code, media)

      getPreviewResp(id, url, method, code, media).then(json => {
        if (json.code === 0) {
          this.respSample = json.data
          this.responseVisible = true
        }
      })
    },

    closePreview() {
      this.responseVisible = false;
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
