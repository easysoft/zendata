<template>
  <a-modal
      title="测试数据设计"
      width="100%"
      dialogClass="full-screen-modal"
      :visible="visible"
      okText="确认"
      cancelText="关闭"
      @ok="save"
      @cancel="cancel"
  >
  DESIGN
  </a-modal>
</template>

<script>

export default {
  name: 'DefDesignComponent',
  data() {
    return {

    };
  },
  props: {
    visible: {
      type: Boolean,
      required: true
    },
    model: {
      type: Object,
      default: () => null
    }
  },

  computed: {
  },
  created () {
  },
  mounted () {

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
  }
}
</script>

<style scoped>

</style>
