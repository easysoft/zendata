<template>
  <div class="mock-preview-main">
    <a-card>
      <span slot="title">
        <a-icon type="profile" />
      </span>

      <div v-if="mockItem">
        <div v-for="(path, url) in mockItem.paths" :key="url" class="path item">
          <div>{{url}}</div>
          <div v-for="(methodVal, method) in path" :key="method" class="method item">
            <div>{{method}}</div>
            <div v-for="(codeVal, code) in methodVal" :key="code" class="code item">
              <div>{{code}}</div>
              <div v-for="(mediaVal, media) in codeVal" :key="media" class="media item">
                <a @click="preview(mockItem.id, url, method, code, media)">{{media}}</a>

                <span :param="fullKey = url+'-'+method+'-'+code+'-'+media">
                  <span :param="samples = dataSrc[fullKey]">
                    <a-select v-if="samples && samples.length > 1"
                              :defaultValue="mockSrcs[fullKey] || samples[0]"
                              @change="selectSample"
                              size="small" class="data-src">
                      <a-select-option v-for="(item, index) in samples" :value="item+'~~~'+fullKey" :key="index">
                        {{item}}
                      </a-select-option>
                    </a-select>
                  </span>
                </span>

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
import {changeSampleSrc, getPreviewResp} from "@/api/mock";

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

    selectSample(value) {
      const arr = value.split('~~~')
      const val = arr[0]
      const key = arr[1]

      changeSampleSrc(this.mockItem.id, key, val)
    },
  },
  watch: {

  }
}
</script>

<style lang="less" scoped>
.mock-preview-main {
  .item {
    line-height: 26px;
  }
  .path {
    padding-left: 10px;
    .method {
      padding-left: 10px;
      .code {
        padding-left: 10px;
        .media {
          padding-left: 10px;
          .data-src {
            margin-left: 10px;
            width: 100px;
          }
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
