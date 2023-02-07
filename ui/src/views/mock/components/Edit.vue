<template>
  <div class="mock-edit-main">
    <a-modal
      :title="$t('msg.mock.create')"
      width="100%"
      dialogClass="full-screen-modal"
      :visible="visible"
      :closable=true
      :footer="null"
      @cancel="cancel"
    >
      <div>

      </div>
    </a-modal>

  </div>
</template>

<script>
import {
} from "../../../api/manage";

export default {
  name: 'MockEditComp',
  components: {

  },
  data() {
    const styl = 'height: ' + (document.documentElement.clientHeight - 56) + 'px;'
    return {
      styl: styl,
      modelData: {},
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
    model: {
      type: Object,
      default: () => {return {}}
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
  mounted: function () {
    console.log('mounted')
  },
  beforeDestroy() {
    console.log('beforeDestroy')
  },

  methods: {
    save() {
      console.log('save')
      this.loadData()
    },
    cancel() {
      console.log('cancel')
      this.$emit('cancel')
    },

    loadData () {
      console.log('loadData', this.modelProp)

      if (!this.model?.id) return

      // getDefFieldTree(this.modelProp.id).then(json => {
      //   console.log('getDefFieldTree', json)
      //   this.loadTreeCallback(json, selectedKey)
      // })
    },

    getModel(id) {
      console.log('getModel', id)
    },
  }
}
</script>

<style lang="less" scoped>
.mock-edit-main {

}

</style>
