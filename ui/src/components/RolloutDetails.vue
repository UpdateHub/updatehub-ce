<template>
<div>
  <div class="pb-0 mt-0 mb-4 border-bottom" v-if="!embedded">
    <h3>
      Rollout Details
      <a href="#" class="btn btn-link float-sm-right" v-if="rollout.running" @click="stop()"><i class="fas fa-ban"></i> Stop</a>
      <router-link to="/rollouts" class="btn btn-link float-sm-right"><i class="fas fa-caret-left"></i> Back to Rollout List</router-link>
    </h3>
  </div>
  <div class="alert d-flex flex-row" v-if="!embedded && ['finised', 'failed'].includes(rollout.statistics.status)" v-bind:class="{ 'alert-danger': rollout.statistics.status == 'failed', 'alert-success': rollout.statistics.status == 'finished'  }" role="alert">
    <i class="fas fa-4x" v-bind:class="{ 'fa-check-circle': rollout.statistics.status == 'finished', 'fa-exclamation-circle': rollout.statistics.status == 'failed' }"></i>
    <div class="align-self-center ml-2">
      <span v-if="rollout.statistics.status == 'finished'">All of <strong>{{ rollout.devices.length }}</strong> devices has been updated successfully!</span>
      <span v-if="rollout.statistics.status == 'failed'"> <strong>{{ rollout.statistics.statuses.failed }}</strong> of <strong>{{ rollout.devices.length }}</strong> devices has been failed while updating!</span>
    </div>
  </div>
  <div class="card-group mb-4" v-if="rollout.package.uid">
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
          <span v-if="rollout.finished_at > rollout.started_at">{{ rollout.finished_at | humanizedDate }}</span>
        </li>
      </ul>
    </div>


    <div class="card card-body col-md-3a">
      <span class="card-title text-center text-capitalize">
        <strong>Status</strong>: {{ rollout.statistics.status }} <span class="text-right"><i v-if="rollout.statistics.status == 'running'" class="fas fa-circle-notch text-mutted" :class="{ 'fa-spin': true}"></i></span>
      </span>
      <ul class="list-group list-group-flush">
        <li class="list-group-item">
          <span><i class="fas fa-question-circle text-mutted"></i> {{ rollout.statistics.statuses.pending }} Pending</span><br/>
        </li>

        <li class="list-group-item">
          <span><i class="fas fa-cog text-primary" :class="{ 'fa-spin': rollout.statistics.statuses.updating > 0 }"></i> {{ rollout.statistics.statuses.updating }} Updating</span><br/>
        </li>

        <li class="list-group-item">
          <span><i class="fas fa-check-circle text-success"></i> {{ rollout.statistics.statuses.updated }} Updated</span><br/>
        </li>
        <li class="list-group-item">
          <span><i class="fas fa-exclamation-circle text-danger"></i> {{ rollout.statistics.statuses.failed }} Failed</span><br/>
        </li>
      </ul>
    </div>
  </div>

 <table class="table table-hover" v-if="!embedded">
        <thead class="thead-light">
          <tr><th>UID</th><th>Version</th><th>Hardware</th><th>Status</th><th></th></tr>
        </thead>
        <tbody>
            <template v-for="device in rollout.devices">
            <tr :key="device.uid" @click="toggleOpened(device.uid)" v-bind:class="deviceRowContextualClass(device)">
              <td><span>{{ device.uid }}</span></td>
              <td><span>{{ device.version }}</span></td>
              <td><span>{{ device.hardware }}</span></td>
              <td><span class="text-capitalize">{{ device.status }}</span></td>
              <td><i class="fas" v-bind:class="{ 'fa-caret-up': opened == device.uid, 'fa-caret-down': opened != device.uid }"></i></td>
            </tr>
            <tr :key="device.uid + 'details'" v-if="opened == device.uid">
              <td colspan="5">
                <div class="alert alert-danger">
                  {{ lastRolloutReportForDevice.message }}
                </div>
                <DeviceDetails :uid="device.uid" embedded="true"></DeviceDetails>
              </td>
            </tr>
            </template>
        </tbody>
      </table>
</div>
</template>

<script>
import DeviceDetails from "./DeviceDetails";

export default {
  name: "RolloutDetails",

  components: { DeviceDetails },

  props: ["id", "embedded"],

  data() {
    return {
      rollout: { package: {}, devices: [], statistics: {}, status: "" },
      timer: null,
      opened: "",
      lastRolloutReportForDevice: {}
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
      const id = this.$route.params.id || this.id;
      this.rollout = await this.getRollout(id).then(async rollout => {
        rollout.package = await this.getPackage(rollout.package);
        rollout.statistics = await this.getStatistics(rollout);
        rollout.devices = await this.getDevices(rollout);

        return rollout;
      });
    },

    async getRollout(id) {
      return await this.$http.get("/api/rollouts/" + id).then(res => {
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
    },

    async getDevices(rollout) {
      return await this.$http
        .get("/api/rollouts/" + rollout.id + "/devices")
        .then(res => {
          return res.data;
        });
    },

    stop() {
      this.$http.put("/api/rollouts/" + this.rollout.id + "/stop").then(res => {
        this.refresh();
      });
    },

    toggleOpened(uid) {
      if (this.opened == uid) {
        this.opened = "";
        this.lastRolloutReportForDevice = {};
      } else {
        this.opened = "";
        this.$nextTick(() => {
          this.opened = uid;
        });

        this.$http
          .get(
            "/api/devices/" + uid + "/rollouts/" + this.rollout.id + "/reports"
          )
          .then(res => {
            this.lastRolloutReportForDevice = res.data[res.data.length - 1];
          });
      }
    },

    deviceRowContextualClass(device) {
      return {
        updated: "table-success",
        failed: "table-danger",
        pending: "table-secondary"
      }[device.status];
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