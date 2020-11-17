<template>
  <div id="design-page">
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
            @rightClick="onRightClick"
            :draggable="true"
            @dragenter="onDragEnter"
            @drop="onDrop"
        />
        <div v-if="treeNode" :style="this.tmpStyle" class="tree-context-menu">
          <a-menu @click="menuClick" mode="inline" class="menu">
            <a-menu-item key="addNeighbor" v-if="!isRoot">
              <a-icon type="plus" />创建同级
            </a-menu-item>
            <a-menu-item key="addChild">
              <a-icon type="plus" />创建子级
            </a-menu-item>
            <a-menu-item key="remove" v-if="!isRoot">
              <a-icon type="delete" />删除节点
            </a-menu-item>
          </a-menu>
        </div>
      </div>

      <div class="right" :style="styl">
        <div v-if="rightVisible">
          <a-tabs :activeKey="tabKey" @change="onChange" type="card">
          <a-tab-pane key="info" tab="编辑信息">
            <div>
              <field-info-component
                  ref="infoComp"
                  :model="fieldModel"
                  @save="onFieldSave">
              </field-info-component>
            </div>
          </a-tab-pane>

          <a-tab-pane key="range" tab="配置区间" force-render>
            <div>
              <field-range-component
                  :field="fieldModel">
              </field-range-component>
            </div>
          </a-tab-pane>

          <a-tab-pane key="refer" tab="配置引用" force-render>
            <div>
              <field-refer-component
                  :field="fieldModel">
              </field-refer-component>
            </div>
          </a-tab-pane>

        </a-tabs>
        </div>
      </div>
    </div>
    </a-modal>

    <a-modal
        title="确认删除"
        :width="400"
        :visible="removeVisible"
        okText="确认"
        cancelText="取消"
        @ok="removeField"
        @cancel="cancelRemove">
      <div>确认删除选中字段及其子字段？</div>
    </a-modal>

  </div>
</template>

<script>
import { getDefFieldTree, getDefField, createDefField, removeDefField, moveDefField } from "../api/manage";
import FieldInfoComponent from "./FieldInfo";
import FieldRangeComponent from "./FieldRange";
import FieldReferComponent from "./FieldRefer";

export default {
  name: 'DefDesignComponent',
  components: {
    FieldInfoComponent,
    FieldRangeComponent,
    FieldReferComponent
  },
  data() {
    const styl = 'height: ' + (document.documentElement.clientHeight - 56) + 'px;'
    return {
      styl: styl,
      removeVisible: false,

      tabKey: 'info',
      rightVisible: true,
      fieldModel: {},
      time2: 0,

      treeData: [],
      openKeys: [],
      selectedKeys: [],
      targetModel: 0,
      treeNode: null,
      fieldMap: {title:'field', key:'id', value: 'id'},
    };
  },
  props: {
    type: {
      type: Number,
      required: true
    },
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
    isRoot () {
      console.log('isRoot', !this.treeNode.parentID)
      return !this.treeNode.parentID
    },
  },
  created () {
    console.log('created')
    this.loadTreeData()
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadTreeData()
    })
  },
  mounted: function () {
    console.log('mounted')
    window.addEventListener("click", this.clearMenu)
  },
  beforeDestroy() {
    console.log('beforeDestroy')
    window.removeEventListener('click', this.clearMenu);
  },
  methods: {
    onFieldSave() {
      console.log('onFieldSave')
      this.loadTreeData(this.selectedKeys[0])
    },
    cancel() {
      console.log('cancel')
      this.$emit('cancel')
    },

    loadTreeData (selectedKey) {
      if (!this.model.id)
        return

      getDefFieldTree(this.model.id).then(res => {
        console.log('getDefFieldTree', res)
        if (res.code != 1) return
        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        if (selectedKey) {
          this.getField(selectedKey)
          this.rightVisible = true
        } else {
          this.rightVisible = false
        }
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
      if (selectedKeys.length == 0) {
        selectedKeys[0] = e.node.eventKey // keep selected
      }

      this.getField(parseInt(selectedKeys[0]))
    },
    getField(id) {
      getDefField(id).then(res => {
        console.log('getDefField', res)
        this.fieldModel = res.data
        this.time2 = Date.now() // trigger data refresh

        this.tabKey = 'info'
        if (this.fieldModel.parentID == 0) {
          this.rightVisible = false
        } else {
          this.rightVisible = true
        }
      })
    },
    menuClick (e) {
      console.log('menuClick', e, this.treeNode)
      this.addMode = null

      this.targetModel = this.treeNode.id
      if (e.key === 'addNeighbor') {
        this.addMode = 'neighbor'
        this.addNeighborField()
      } else if (e.key === 'addChild') {
        this.addMode = 'child'
        this.addChildField()
      }else if (e.key === 'remove') {
        this.removeVisible = true
      }
      console.log('clearMenu 1')
      this.clearMenu()
    },
    addNeighborField () {
      console.log('addNeighborField', this.targetModel)

      createDefField(this.targetModel, "neighbor").then(res => {
        console.log('createDefField', res)

        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        this.selectedKeys = [res.field.id] // select
        this.fieldModel = res.field

        this.rightVisible = true
      })
    },
    addChildField () {
      console.log('addChildField', this.targetModel)

      createDefField(this.targetModel, "child").then(res => {
        console.log('createDefField', res)

        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        this.selectedKeys = [res.field.id] // select
        this.fieldModel = res.field

        this.rightVisible = true
      })
    },
    removeField () {
      console.log('removeField', this.targetModel)
      this.removeVisible = false

      removeDefField(this.targetModel).then(res => {
        console.log('removeDefField', res)

        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        this.rightVisible = false
      })
    },
    cancelRemove (e) {
      e.preventDefault()
      this.removeVisible = false
    },
    onDragEnter(info) {
      console.log(info);
      // expandedKeys 需要受控时设置
      this.expandedKeys = info.expandedKeys
    },
    onDrop(info) {
      console.log(info, info.dragNode.eventKey, info.node.eventKey, info.dropPosition);

      moveDefField(info.dragNode.eventKey, info.node.eventKey, info.dropPosition).then(res => {
        this.getOpenKeys(res.data)
        this.treeData = [res.data]

        this.selectedKeys = [res.field.id] // select
        this.fieldModel = res.field

        this.rightVisible = true
      })
    },

    onRightClick ({ event, node }) {
      event.preventDefault()
      console.log('onRightClick', node)

      const y = event.currentTarget.getBoundingClientRect().top
      const x = event.currentTarget.getBoundingClientRect().right

      this.treeNode = {
        pageX: x,
        pageY: y,
        id: node._props.eventKey,
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
      console.log('clearMenu')
      this.treeNode = null
    },
    onChange(activeKey) {
      console.log('onChange', activeKey)
      this.tabKey = activeKey
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

.tree-context-menu {
  z-index: 9;
  .ant-tree-node-content-wrapper {
    display: block !important;
  }
  .menu {
    border: 1px solid #ebedf0;
    background: #f0f2f5;
    .ant-menu-item {
      padding-left: 12px !important;
      height: 22px;
      line-height: 21px;
    }
  }
}

</style>
