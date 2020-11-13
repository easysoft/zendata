<template>
  <div class="panel">
    <div class="radios">
      <a-radio-group v-model="model.isRange" button-style="solid">
        <a-radio-button :value="true">
          区间
        </a-radio-button>
        <a-radio-button :value="false">
          引用
        </a-radio-button>
      </a-radio-group>
      &nbsp;&nbsp;&nbsp;
      <span class="range">{{model.range}}</span>
    </div>

    <div v-if="model.isRange">

    </div>

    <div v-if="!model.isRange">

    </div>

  </div>
</template>

<script>

import FieldRangeComponent from "./FieldRange";
import FieldReferComponent from "./FieldRefer";

export default {
  name: 'FieldConfigComponent',
  components: {FieldRangeComponent, FieldReferComponent},
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
  },
  mounted () {
    console.log('mounted')
  },
  methods: {

  }
}
</script>

<style lang="less" scoped>

</style>
