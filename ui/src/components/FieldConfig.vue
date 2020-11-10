<template>
  <div>
    CONFIG
  </div>
</template>

<script>
import { getDefField } from "../api/manage";

export default {
  name: 'FieldConfigComponent',
  data() {
    return {
    };
  },
  props: {
    model: {
      type: Object,
      default: () => null
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
  mounted () {
    console.log('mounted1')
  },
  methods: {
    save() {
      console.log('save')
      this.$emit('ok')
    },
    cancel() {
      console.log('cancel')
      this.$emit('cancel')
    },

    loadData () {
      if (!this.model.id) return

      getDefField(this.model.id).then(res => {
        console.log('getField', res)
        this.model = [res.data]
      })
    },
  }
}
</script>

<style lang="less" scoped>

</style>
