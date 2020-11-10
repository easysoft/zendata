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
          :show-line="true"
          :expandedKeys.sync="openKeys"
          :selectedKeys.sync="selectedKeys"
          :tree-data="treeData"
          :replaceFields="fieldMap"
          @select="onSelect"
          @rightClick="this.onRightClick"
      />
    </div>

    <div class="right" :style="styl">
      <a-tabs default-active-key="1" @change="onChange">
        <a-tab-pane key="info" tab="编辑">
          <div v-show="infoVisible">
            <field-info
                ref="infoComp"
                :model="fieldModel"
                :time="time2">
            </field-info>
          </div>
        </a-tab-pane>

        <a-tab-pane key="config" tab="设计" force-render>
          <div v-show="configVisible">
            <field-config
                ref="configComp"
                :model="fieldModel"
                :time="time2">
            </field-config>
          </div>
        </a-tab-pane>
      </a-tabs>
    </div>
  </div>
  </a-modal>
</template>

<script>
import { getDefFieldTree, getDefField } from "../api/manage";
import FieldInfoComponent from "./FieldInfo";
import FieldConfigComponent from "./FieldConfig";

export default {
  name: 'DefDesignComponent',
  components: {
    'field-info': FieldInfoComponent,
    'field-config': FieldConfigComponent
  },
  data() {
    const styl = 'height: ' + (document.documentElement.clientHeight - 56) + 'px;'
    return {
      styl: styl,

      infoVisible: true,
      configVisible: false,
      fieldModel: {},
      time2: 0,

      treeData: [],
      openKeys: [],
      selectedKeys: [],
      treeNode: null,
      fieldMap: {title:'field', key:'id', value: 'id'},
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
    console.log('created')
    this.loadTreeData()
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadTreeData()
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

    loadTreeData () {
      if (!this.model.id) return

      getDefFieldTree(this.model.id).then(res => {
        console.log('getDefFieldTree', res)
        if (res.code != 1) return
        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        this.infoVisible = false
        this.configVisible = false
      })
    },
    getOpenKeys (def) {
      if (!def) return

      this.openKeys.push(def.id)
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

      getDefField(parseInt(selectedKeys[0])).then(res => {
        console.log('getDefField', res)
        this.fieldModel = res.data

        if (this.fieldModel.parentID == 0) {
          this.infoVisible = false
          this.configVisible = false
        } else {
          this.infoVisible = true
          this.configVisible = true
        }
      })
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
    },
    onChange() {
      console.log('onChange')
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
    overflow: auto;
  }
  .right {
    flex: 1;
    height: 100%;
    padding: 6px;
    overflow: auto;
  }
}

</style>
