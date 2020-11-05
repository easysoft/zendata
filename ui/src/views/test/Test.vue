<template>
  <div class="test">
    <div>
      TEST
    </div>
    <div>
      {{$t('common.submit')}}
    </div>
    <a-button type="primary">Button></a-button>
  </div>
</template>

<script>

import { listDefs } from "../../api/manage";

export default {
  name: 'Test',
  data() {
    return {
      defs: [],
    };
  },
  computed: {

  },
  created () {
    listDefs().then(res => {
      console.log('listDefs', res)
      this.defs = res
    })
  },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped>

</style>
