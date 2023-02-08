import { mapState } from 'vuex'

const mockMixin = {
  computed: {
    ...mapState({
      mockItem: state => state.mock.mockItem
    })
  },
  methods: {
    previewMockItem (item) {
      this.$store.dispatch('previewMockItem', item.id)
    },
    saveMockItem (item) {
      return this.$store.dispatch('saveMockItem', item)
    },
  }
}

export default mockMixin
