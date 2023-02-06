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
    }
  }
}

export default mockMixin
