<template>
  <div class="yaml-editor">
    <textarea ref="textarea"/>
  </div>
</template>

<script>
import CodeMirror from 'codemirror'
import 'codemirror/addon/lint/lint.css'
import 'codemirror/lib/codemirror.css'
import 'codemirror/mode/yaml/yaml'
import 'codemirror/addon/lint/lint'
import 'codemirror/addon/lint/yaml-lint'
// 提示弹窗
import 'codemirror/addon/dialog/dialog.js'
import 'codemirror/addon/dialog/dialog.css'
// 滚动条
import 'codemirror/addon/scroll/simplescrollbars.css'
import 'codemirror/addon/scroll/simplescrollbars.js'
// 搜索功能
import 'codemirror/addon/search/search.js'
import 'codemirror/addon/search/searchcursor.js'
import 'codemirror/addon/search/jump-to-line.js'
// 代码高亮
import "codemirror/addon/selection/active-line";

window.jsyaml = require("js-yaml") // 引入js-yaml为codemirror提高语法检查核心支持

export default {
  name: 'YamlEditor',
  props: ['value'],
  data() {
    return {
      yamlEditor: false
    }
  },
  watch: {
    value(value) {
      const editorValue = this.yamlEditor.getValue()
      if (value !== editorValue) {
        this.yamlEditor.setValue(this.value)
      }
    }
  },
  mounted() {
    this.yamlEditor = CodeMirror.fromTextArea(this.$refs.textarea, {
      lineNumbers: true, // 显示行号
      mode: 'text/x-yaml', // 语法model
      gutters: ['CodeMirror-lint-markers'],  // 语法检查器
      lint: true, // 开启语法检查
      indentUnit: 1,         // 缩进单位为2
      styleActiveLine: true, // 当前行背景高亮
      matchBrackets: true,   // 括号匹配
      lineWrapping: true,    // 自动换行
      tabSize: 2,
      smartIndent: true,
    })

    this.yamlEditor.setValue(this.value)
    this.yamlEditor.on('change', (cm) => {
      this.$emit('changed', cm.getValue())
      this.$emit('input', cm.getValue())
    })
  },
  methods: {
    getValue() {
      return this.yamlEditor.getValue()
    }
  }
}
</script>

<style scoped>
.yaml-editor {
  height: 100%;
  position: relative;
}

.yaml-editor >>> .CodeMirror {
  height: auto;
  min-height: 300px;
}

.yaml-editor >>> .CodeMirror-scroll {
  min-height: 300px;
}

.yaml-editor >>> .cm-s-rubyblue span.cm-string {
  color: #F08047;
}
</style>