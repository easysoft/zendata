<template>
  <div id="design-page">
    <div class="container">
    <div class="left" :style="styl">
    <a-tree
        ref="fieldTree"
        class="draggable-tree"
        :show-line="true"
        :show-icon="false"
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
            <a-icon type="plus" />{{ $t('msg.design.create.brother') }}
        </a-menu-item>
        <a-menu-item key="addChild" v-if="type=='def'|| ((type=='ranges' || type=='instances') && isRoot)">
            <a-icon type="plus" />{{ $t('msg.design.create.child') }}
        </a-menu-item>
        <a-menu-item key="remove" v-if="!isRoot">
            <a-icon type="delete" />{{ $t('msg.design.remove.node') }}
        </a-menu-item>
        </a-menu>
    </div>
    </div>

    <div class="right" :style="styl">
    <div v-if="rightVisible">

        <div v-if="type=='def' || type=='instances'">
        <a-tabs :activeKey="tabKey" @change="onTabChange" type="card">
            <a-tab-pane key="info" :tab="$t('msg.info')">
            <div>
                <field-info-component
                    ref="infoComp"
                    :type="type"
                    :model="modelData"
                    @save="onModelSave">
                </field-info-component>
            </div>
            </a-tab-pane>

            <a-tab-pane key="range" :tab="$t('msg.range')" force-render>
            <div>
                <field-range-component
                    ref="rangeComp"
                    :type="type"
                    :model="modelData"
                    :time2="time2">
                </field-range-component>
            </div>
            </a-tab-pane>

            <a-tab-pane key="refer" :tab="$t('msg.reference')" force-render>
            <div>
                <field-refer-component
                    ref="referComp"
                    :type="type"
                    :model="modelData"
                    :time2="time2">
                </field-refer-component>
            </div>
            </a-tab-pane>

            <a-tab-pane key="preview" :tab="$t('msg.preview')" force-render>
            <pre class="preview-data" v-html="previewData" style="margina: 0"></pre>
            </a-tab-pane>
        </a-tabs>
        </div>

        <div v-if="type=='ranges'">
        <res-ranges-item-component
            ref="rangesItem"
            :model="modelData"
            :time="time2"
            @save="onModelSave">
        </res-ranges-item-component>
        </div>

        <div v-if="type=='config'"> <!-- no item object, show sections page -->
        <div class="head">
            <div class="title">
            字段编辑
            </div>
            <div class="buttons"></div>
        </div>

        <a-row>
            <a-col :offset="2">
            <field-range-component
                ref="rangeComp"
                :type="'config'"
                :model="modelData"
                :time2="time2">
            </field-range-component>
            </a-col>
        </a-row>
        </div>

    </div>
    </div>
</div>

    <a-modal
        title="确认删除"
        :width="400"
        :visible="removeVisible"
        okText="确认"
        cancelText="取消"
        @ok="removeNode"
        @cancel="cancelRemove">
      <div>确认删除选中节点？</div>
    </a-modal>

  </div>
</template>

<script>
import { getDefFieldTree, getDefField, createDefField, removeDefField, moveDefField,
         getResRangesItemTree, getResRangesItem, createResRangesItem, removeResRangesItem,
         getResInstancesItemTree, getResInstancesItem, createResInstancesItem, removeResInstancesItem,
         getResConfigItemTree,
} from "../api/manage";
import FieldInfoComponent from "./FieldInfo";
import FieldRangeComponent from "./FieldRange";
import FieldReferComponent from "./FieldRefer";
import ResRangesItemComponent from "./RangesItem"
import {ResTypeDef, ResTypeInstances, ResTypeRanges, ResTypeConfig} from "../api/utils";
import {previewFieldData} from "../api/manage";

export default {
  name: 'DefDesignComponent',
  components: {
    FieldInfoComponent, FieldRangeComponent, FieldReferComponent,
    ResRangesItemComponent,
  },
  data() {
    const styl = 'height: ' + (document.documentElement.clientHeight - 56) + 'px;'
    return {
      styl: styl,
      removeVisible: false,

      tabKey: 'info',
      rightVisible: true,
      modelData: {},
      time2: 0,

      previewData: '',
      treeData: [],
      nodeMap: {},
      openKeys: [],
      selectedKeys: [],
      targetModel: 0,
      treeNode: null,
      fieldMap: {title: 'field', key:'id', value: 'id', children: 'fields'},
    };
  },
  props: {
    type: {
      type: String,
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
    isRoot () {
      console.log('isRoot', this.treeNode)
      return !this.treeNode.parentID || this.treeNode.parentID == 0 || this.treeNode.id == 0
    },
  },
  created () {
    console.log('created')
    this.loadTree()
    this.$watch('time', () => {
      console.log('time changed', this.time)
      this.loadTree()
    })
  },
  mounted: function () {
    console.log('mounted')
  },
  beforeDestroy() {
    console.log('beforeDestroy')
  },

  watch: {
    modelProp(val) {
      console.log("watch modelProp :", val)
      this.type = ResTypeDef;
      this.loadTree('')
    },
  },

  methods: {
    onModelSave() {
      console.log('onModelSave')
      this.$emit('save')
      this.loadTree(this.selectedKeys[0])
    },
    cancel() {
      console.log('cancel')
      this.$emit('cancel')
    },

    loadTree (selectedKey) {
      console.log('loadTree', this.modelProp)
      if (!this.modelProp.id)
        return

      if (this.type === ResTypeDef) {
        getDefFieldTree(this.modelProp.id).then(json => {
          console.log('getDefFieldTree', json)
          this.loadTreeCallback(json, selectedKey)
        })
      } else if (this.type === ResTypeRanges) {
        getResRangesItemTree(this.modelProp.id).then(json => {
          console.log('getResRangesItemTree', json)
          this.loadTreeCallback(json, selectedKey)
        })
      } else if (this.type === ResTypeInstances) {
        getResInstancesItemTree(this.modelProp.id).then(json => {
          console.log('getResInstancesItemTree', json)
          this.loadTreeCallback(json, selectedKey)
        })
      } else if (this.type === ResTypeConfig) {
        getResConfigItemTree(this.modelProp.id).then(json => {
          console.log('getResConfigItemTree', json)
          this.selectedKeys = [this.modelProp.id]
          this.loadTreeCallback(json, this.modelProp.id)
        })
      }
    },
    loadTreeCallback(json, selectedKey) {
      if (json.code != 1) return
      this.getOpenKeys(json.data)
      this.treeData = [json.data]

      if (selectedKey) {
        this.getModel(selectedKey)
        this.rightVisible = true
      } else {
        this.rightVisible = false
      }
    },
    getOpenKeys (node) {
      if (!node) return

      this.openKeys.push(node.id)
      this.nodeMap[node.id] = node
      if (node.fields) {
        node.fields.forEach((item) => {
          this.getOpenKeys(item)
        })
      }
    },
    onSelect (selectedKeys, e) { // selectedKeys, e:{selected: bool, selectedNodes, node, event}
      console.log('onSelect', selectedKeys, e.selectedNodes, e.node, e.node.eventKey)
      if (selectedKeys.length == 0) {
        selectedKeys[0] = e.node.eventKey // keep selected
      }

      const node = this.nodeMap[e.node.eventKey]
      console.log('node', node)
      if ((this.type === 'def' && node.parentID == 0) || (this.type === 'config' && node.id == 0)
          || (node.fields && node.fields.length > 0)) {
        this.rightVisible = false
        this.modelData = {}
        return
      } else {
        this.rightVisible = true
        this.tabKey = 'info'
      }

      this.getModel(parseInt(selectedKeys[0]))
    },
    getModel(id) {
      console.log('getModel', id)

      if (this.type === 'def') {
        getDefField(id).then(res => {
          console.log('getDefField', res)
          this.modelData = res.data
          this.time2 = Date.now() // trigger data refresh
        })
      } else if (this.type === 'ranges') {
        getResRangesItem(id).then(res => {
          console.log('getResRangesItem', res)
          this.modelData = res.data
          this.time2 = Date.now() // trigger data refresh
        })
      } else if (this.type === 'instances') {
        getResInstancesItem(id).then(res => {
          console.log('getResInstancesItem', res)
          this.modelData = res.data
          this.time2 = Date.now() // trigger data refresh
        })
      } else if (this.type === 'config') {
        this.modelData = {id: id}
        this.time2 = Date.now() // trigger data refresh
      }
    },
    menuClick (e) {
      console.log('menuClick', e, this.treeNode)
      this.addMode = null

      this.targetModel = this.treeNode.id
      if (e.key === 'addNeighbor') {
        this.addMode = 'neighbor'
        this.addNeighbor()
      } else if (e.key === 'addChild') {
        this.addMode = 'child'
        this.addChildField()
      }else if (e.key === 'remove') {
        this.removeVisible = true
      }
      this.clearMenu()
    },
    addNeighbor () {
      console.log('addNeighbor', this.targetModel)

      if (this.type === 'def') {
        createDefField(this.targetModel, "neighbor").then(json => {
          console.log('createDefField', json)
          this.updateCallback(json)
        })
      } else if (this.type === 'ranges') {
        createResRangesItem(this.modelProp.id, "neighbor").then(json => {
          console.log('createResRangesItem', json)
          this.updateCallback(json)
        })
      } else if (this.type === 'instances') {
        createResInstancesItem(this.modelProp.id, "neighbor").then(json => {
          console.log('createResInstancesItem', json)
          this.updateCallback(json)
        })
      }
    },
    addChildField () {
      console.log('addChildField', this.targetModel)

      if (this.type === 'def') {
        createDefField(this.targetModel, "child").then(json => {
          console.log('createDefField', json)
          this.updateCallback(json)
        })
      } else if (this.type === 'ranges') {
        createResRangesItem(this.modelProp.id, "child").then(json => {
          console.log('createResRangesItem', json)
          this.updateCallback(json)
        })
      } else if (this.type === 'instances') {
        createResInstancesItem(this.modelProp.id, "child").then(json => {
          console.log('createResInstancesItem', json)
          this.updateCallback(json)
        })
      }
    },
    updateCallback(json) {
      this.getOpenKeys(json.data)
      this.treeData = [json.data]

      this.selectedKeys = [json.model.id] // select
      this.modelData = json.model

      this.rightVisible = true
    },
    removeNode () {
      console.log('removeNode', this.targetModel)
      this.removeVisible = false
      if (this.type === 'def') {
        removeDefField(this.targetModel).then(json => {
          console.log('removeDefField', json)
          this.removeCallback(json)
        })
      } else if (this.type === 'ranges') {
        removeResRangesItem(this.targetModel, this.modelProp.id).then(json => {
          console.log('removeResRangesItem', json)
          this.removeCallback(json)
        })
      } else if (this.type === 'instances') {
        removeResInstancesItem(this.targetModel, this.modelProp.id).then(json => {
          console.log('removeResInstancesItem', json)
          this.removeCallback(json)
        })
      }
    },
    removeCallback(json) {
      this.getOpenKeys(json.data)
      this.treeData = [json.data]

      this.rightVisible = false
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

        this.selectedKeys = [res.model.id] // select
        this.modelData = res.model

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
    onTabChange(activeKey) {
      console.log('onTabChange', activeKey)
      this.tabKey = activeKey

      if (this.tabKey === 'preview') {
        console.log(111, this.getFieldTreeName(this.modelData))
        previewFieldData({config: this.modelProp.referName.replace(/\\/g, "/"), field: this.getFieldTreeName(this.modelData), format: 'txt'}).then(data => {
          console.log('previewFieldData', data)
          this.previewData = data
        })
      }
    },
    getFieldTreeName(node, name = ''){
      if(node.parentID == 0){
        return name;
      }
      if(this.nodeMap[node.parentID] == undefined){
        return name == '' ? node.field : node.field + '~~' + name;
      }

      name = name == '' ? node.field : node.field + '~~' + name;

      return this.getFieldTreeName(this.nodeMap[node.parentID], name)
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
      text-align: left;
    }
  }
}

.preview-data {
  padding: 5px 20px;
}

</style>
