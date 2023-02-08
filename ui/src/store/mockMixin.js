import { mapState } from 'vuex'

const mockMixin = {
  computed: {
    ...mapState({
      mockItem: state => state.mock.mockItem
    })
  },
  methods: {
    setMockItem (item) {
      this.$store.dispatch('setMockItem', item)
    },
    saveMockItem (item) {
      return this.$store.dispatch('saveMockItem', item)
    },
  }
}

export default mockMixin
