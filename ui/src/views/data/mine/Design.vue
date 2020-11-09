<template>
  <div>
    <div class="head">
      <div class="title">
        测试数据设计
      </div>
      <div class="buttons">
        <a-button type="primary" @click="back()">返回</a-button>
      </div>
    </div>

    <div>
      DESIGN
    </div>
  </div>
</template>

<script>

import { getDef, saveDefDesign } from "../../../api/manage";

export default {
  name: 'Mine',
  data() {
    return {
      id: 0,
      model: { folder: '/users' }
    };
  },
  computed: {
  },
  created () {
    this.id = parseInt(this.$route.params.id)
    if (this.id == 0) return

    getDef(this.id).then(json => {
      console.log('getDef', json)
      this.model = json.data
    })
  },
  mounted () {

  },
  methods: {
    saveDesign() {
      saveDefDesign(this.model).then(json => {
        console.log('saveDefDesign', json)
        // this.back()
      })
    },
    back() {
      this.$router.push({path: '/data/mine/index'});
    },
  }
}
</script>

<style scoped>

</style>
