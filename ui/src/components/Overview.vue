<template>
  <div>
    <div class="pb-0 mt-0 mb-4 border-bottom">
      <h3>Running Rollouts</h3>
    </div>
    <RolloutDetails v-for="rollout in rollouts" :id="rollout.id" :key="rollout.id" embedded="true"/>
    <div class="alert d-flex flex-row" v-if="rollouts.length == 0">
      <div class="col text-center">
        <i class="fas fa-sync fa-spin fa-6x"></i>
        <div class="align-self-center ml-2 mt-4">Currently there is no any active rollout</div>
      </div>
    </div>
  </div>
</template>

<script>
import RolloutDetails from './RolloutDetails'

export default {
  name: 'Overview',

  components: { RolloutDetails },

  data () {
    return {
      rollouts: [],
      timer: null
    }
  },

  async created () {
    this.refresh()
    this.timer = setInterval(this.refresh, 5 * 1000)
  },

  methods: {
    async refresh () {
      this.rollouts = await this.getRollouts()
      this.rollouts = this.filterBy(this.rollouts, { running: true })
    },

    async getRollouts () {
      return this.$http.get('/api/rollouts').then(res => {
        return res.data
      })
    },

    filterBy (values, filter) {
      return values.filter(value => {
        return Object.keys(filter).some(key => {
          return value[key] === filter[key]
        })
      })
    }
  }
}
</script>

<style scoped>
.alert {
  padding: 0;
}

.fa-sync {
  color: #e5e5e5;
}
</style>
