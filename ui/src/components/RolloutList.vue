<template>
<div>
  <div class="pb-0 mt-0 mb-4 border-bottom">
    <h3>
      Rollout List
      <router-link to="/rollouts/new" class="btn btn-link float-sm-right">
        <i class="fas fa-plus"></i> Create Rollout
      </router-link>
    </h3>
  </div>

  <table class="table table-bordered table-hover">
    <thead class="thead-light">
      <tr>
        <th>Version</th><th>Number of Devices</th><th>Created</th><th>Running</th>
      </tr>
    </thead>
    <tbody>
        <tr v-for="rollout in rollouts" :key="rollout.id" @click="$router.push(`/rollouts/${rollout.id}`)">
          <td>{{ rollout.package.version }}</td>
          <td>{{ rollout.devices.length }}</td>
          <td>{{ rollout.started_at | humanizeDate }}</td>
          <td>{{ rollout.running }}</td>
        </tr>
    </tbody>
  </table>
</div>

</template>

<script>
export default {
  name: "RolloutList",

  data() {
    return {
      rollouts: []
    };
  },

  async created() {
    this.rollouts = await this.getRollouts();
    this.rollouts.forEach(async rollout => {
      rollout.package = await this.getPackage(rollout.package);
    });
  },

  methods: {
    async getRollouts() {
      return await this.$http.get("/api/rollouts").then(res => {
        return res.data;
      });
    },

    async getPackage(uid) {
      return await this.$http.get("/api/packages/" + uid).then(res => {
        return res.data;
      });
    }
  },

  filters: {
    humanizeDate(v) {
      return moment(v).format("lll");
    }
  }
};
</script>

<style scoped>
tr {
  cursor: pointer;
}
</style>