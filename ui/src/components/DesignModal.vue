<template>
  <div id="design-page">
    <a-modal
      :title="$t('msg.design.title')"
      width="100%"
      dialogClass="full-screen-modal"
      :visible="visible"
      :closable=true
      :footer="null"
      @cancel="cancel"
    >
    <design-component
    ref="designPage"
    :type="type"
    :modelProp="modelProp"
    :time="time" >
    ></design-component>
    </a-modal>
  </div>
</template>

<script>
import DesignComponent from "./Design"

export default {
  name: 'DefDesignModalComponent',
  components: {
    DesignComponent,
  },
  data() {
    return {
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
    modelProp: {
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
    this.$watch('visible', () => {
      console.log('visible changed', this.visible)
      if (this.visible) {
        document.addEventListener("click", this.clearMenu)
      } else {
        document.removeEventListener('click', this.clearMenu);
      }
    })
  },
  mounted: function () {
    console.log('mounted')
  },
  beforeDestroy() {
    console.log('beforeDestroy')
  },
  methods: {
    cancel() {
      console.log('cancel')
      this.$emit('cancel')
    },
  }
}
</script>

<style lang="less" scoped>
</style>
