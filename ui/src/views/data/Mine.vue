<template>
  <div>

  </div>
</template>

<script>

import { listDef } from "../../api/manage";

export default {
  name: 'Mine',
  data() {
    return {
      defs: [],
    };
  },
  computed: {

  },
  created () {
    console.log('===')
    listDef().then(res => {
      console.log('listDefs', res)
      this.defs = res.data
    })

    // const def = {name: "myDef"}
    // saveDef(def).then(res => {
    //   console.log('saveDef', res)
    //   this.defs = res
    // })
  },
  mounted () {
  },
  methods: {

  }
}
</script>

<style scoped>

</style>
