<template>
<div>
  <div class="pb-0 mt-0 mb-4 border-bottom">
    <h3>
      Rollout Details
      <router-link to="/packages" class="btn btn-link float-sm-right" v-if="rollout.running"><i class="fas fa-pause"></i> Pause</router-link>
      <router-link to="/packages" class="btn btn-link float-sm-right" v-if="!rollout.running"><i class="fas fa-play"></i> Play</router-link>
      <router-link to="/rollouts" class="btn btn-link float-sm-right"><i class="fas fa-caret-left"></i> Back to Rollout List</router-link>
    </h3>
  </div>
  <div class="card-group" v-if="rollout.package.uid">
    <div class="card card-body">
      <ul class="list-group list-group-flush">
        <li class="list-group-item">
          <span>Version</span><br/>
          {{ rollout.package.version }}
        </li>
        <li class="list-group-item">
          <span>Number of Devices</span><br/>
          {{ rollout.devices.length }}
        </li>
        <li class="list-group-item">
          <span>Package</span><br/>
          <router-link :to="`/packages/${rollout.package.uid}`">{{ rollout.package.uid.substring(0, 6) }}...</router-link>
        </li>
      </ul>
    </div>

    <div class="card card-body">
      <ul class="list-group list-group-flush">
        <li class="list-group-item">
          <span>Started At</span><br/>
          {{ rollout.started_at | humanizedDate }}
        </li>
        <li class="list-group-item">
          <span>Finished At</span><br/>
          {{ rollout.started_at > rollout.finished_at ? '-' : rollout.finished_at | humanizedDate }}
        </li>
      </ul>
    </div>


    <div class="card card-body col-md-3a">
      <span class="card-title text-center text-capitalize">
        <strong>Status</strong>: {{ rollout.status }} <span class="text-right"><i v-if="rollout.status == 'running'" class="fas fa-circle-notch text-mutted" :class="{ 'fa-spin': true}"></i></span>
      </span>
      <ul class="list-group list-group-flush">
        <li class="list-group-item">
          <span><i class="fas fa-question-circle text-mutted"></i> {{ rollout.statistics.pending }} Pending</span><br/>
        </li>

        <li class="list-group-item">
          <span><i class="fas fa-cog text-primary" :class="{ 'fa-spin': rollout.statistics.updating > 0 }"></i> {{ rollout.statistics.updating }} Updating</span><br/>
        </li>

        <li class="list-group-item">
          <span><i class="fas fa-check-circle text-success"></i> {{ rollout.statistics.updated }} Updated</span><br/>
        </li>
        <li class="list-group-item">
          <span><i class="fas fa-exclamation-circle text-danger"></i> {{ rollout.statistics.failed }} Failed</span><br/>
        </li>
      </ul>
    </div>
  </div>
</div>
</template>

<script>
export default {
  name: "RolloutDetails",

  data() {
    return {
      rollout: { package: {}, devices: [], statistics: {}, status: "" },
      timer: null
    };
  },

  async created() {
    this.refresh();

    this.timer = setInterval(this.refresh, 5 * 1000);
  },

  beforeDestroy() {
    clearInterval(this.timer);
  },

  methods: {
    async refresh() {
      this.rollout = await this.getRollout().then(async rollout => {
        rollout.package = await this.getPackage(rollout.package);
        rollout.statistics = await this.getStatistics(rollout);

        if (rollout.finished_at > rollout.started_at) {
          rollout.status = "finished";
        } else if (rollout.running) {
          rollout.status = "running";
        } else {
          rollout.status = "paused";
        }

        return rollout;
      });
    },

    async getRollout() {
      return await this.$http
        .get("/api/rollouts/" + this.$route.params.id)
        .then(res => {
          return res.data;
        });
    },

    async getPackage(uid) {
      return await this.$http.get("/api/packages/" + uid).then(res => {
        return res.data;
      });
    },

    async getStatistics(rollout) {
      return await this.$http
        .get("/api/rollouts/" + rollout.id + "/statistics")
        .then(res => {
          return res.data;
        });
    }
  },

  filters: {
    pretty: function(value) {
      return JSON.stringify(value, null, 2)
        .replace(/\n/g, "<br/>")
        .replace(/ /g, "&nbsp;");
    },

    humanizedDate(v) {
      return moment(v).format("lll");
    }
  }
};
</script>
