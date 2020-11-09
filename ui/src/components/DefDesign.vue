<template>
  <a-modal
      title="测试数据设计"
      width="100%"
      dialogClass="full-screen-modal"
      :visible="visible"
      :closable=true
      :footer="null"
      @cancel="cancel"
  >
  <div class="container">
    <div class="left" :style="styl">
      <a-tree
          ref="fieldTree"
          class="draggable-tree"
          :expandedKeys.sync="openKeys"
          :selectedKeys.sync="selectedKeys"
          :tree-data="treeData"
          @select="onSelect"
          @rightClick="this.onRightClick"
      />
    </div>

    <div class="right" :style="styl">right</div>
  </div>
  </a-modal>
</template>

<script>
import { getDefFieldTree } from "../api/manage";

export default {
  name: 'DefDesignComponent',
  data() {
    const styl = 'height: ' + (document.documentElement.clientHeight - 56) + 'px;'
    return {
      styl: styl,

      treeData: [],
      openKeys: [],
      selectedKeys: [],
      treeNode: null,
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
    },
    time: {
      type: Number,
      default: () => 0
    },
  },

  computed: {
  },
  created () {
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadTreeData()
    })
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

    loadTreeData () {
      if (!this.model.id) return

      getDefFieldTree(this.model.id).then(res => {
        console.log('getDefFieldTree', res)
        if (res.code != 1) return
        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        // if (this.selectedKeys.length === 0) {
        //   this.selectedKeys[0] = res.data.key // select root
        //   this.selectKey = res.data.key
        //   // this.$refs.fields.refresh()
        // }
      })
    },
    getOpenKeys (def) {
      if (!def) return

      this.openKeys.push(def.key)
      if (def.children) {
        def.children.forEach((item) => {
          this.getOpenKeys(item)
        })
      }
    },
    onSelect (selectedKeys, e) { // selectedKeys, e:{selected: bool, selectedNodes, node, event}
      console.log('onSelect', selectedKeys, e.node.eventKey)
      if (selectedKeys.length > 0) {
        this.selectKey = selectedKeys[0]
      } else {
        selectedKeys[0] = e.node.eventKey // keep selected
      }

      this.$refs.usersTable.refresh()
    },
    onRightClick ({ event, node }) {
      event.preventDefault()
      const y = event.currentTarget.getBoundingClientRect().top
      const x = event.currentTarget.getBoundingClientRect().right

      console.log('onRightClick', node)
      this.treeNode = {
        pageX: x,
        pageY: y,
        orgID: node._props.eventKey,
        title: node._props.title,
        parentID: node._props.dataRef.parentID || null
      }

      this.tmpStyle = {
        position: 'fixed',
        maxHeight: 40,
        textAlign: 'center',
        left: `${x + 10 - 0}px`,
        top: `${y + 6 - 0}px`
        // display: 'flex',
        // flexDirection: 'row'
      }
    },
    clearMenu () {
      this.treeNode = null
    }
  }
}
</script>

<style lang="less" scoped>
.container {
  display: flex;

  .left {
    padding: 6px;
    width: 220px;
    height: 100%;
    border-right: 1px solid #e9f2fb;
  }
  .right {
    flex: 1;
    height: 100%;
    padding: 6px;
  }
}

</style>
