<template>
  <div class="mock-preview-main">
    <a-card>
      <span slot="title">
        <a-icon type="profile" />
      </span>

      <div>
        {{this.mockItem}}
      </div>

    </a-card>
  </div>
</template>

<script>

import mockMixin from "@/store/mockMixin";

export default {
  name: 'Preview',
  components: {
  },
  props: {
  },
  mixins: [mockMixin],
  data: function() {
    return {};
  },
  mounted: function() {
  },
  methods: {

  },
  watch: {

  }
}
</script>

<style lang="less" scoped>
.mock-preview-main {
}
</style>
