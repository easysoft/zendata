import { mapState } from 'vuex'

const mockMixin = {
  computed: {
    ...mapState({
      mockItem: state => state.mock.mockItem,
      mockSrcs: state => state.mock.mockSrcs,
      dataSrc: state => state.mock.dataSrc,
    })
  },
  methods: {
    setMockItem (item) {
      this.$store.dispatch('setMockItem', item)
    },
    previewMockItem (item) {
      this.$store.dispatch('previewMockItem', item)
    },
    saveMockItem (item) {
      return this.$store.dispatch('saveMockItem', item)
    },
  }
}

export default mockMixin
