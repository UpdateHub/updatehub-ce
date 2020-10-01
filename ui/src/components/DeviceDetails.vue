<template>
<div>
  <template v-if="!embedded">
  <div class="pb-0 mt-0 mb-4 border-bottom">
    <h3>
      Device Details
      <router-link to="/devices" class="btn btn-link float-sm-right"><i class="fas fa-caret-left"></i> Back to Device List</router-link>
    </h3>
  </div>

  <table class="table">
    <thead class="thead-light">
      <tr><th>UID</th><th class="text-center">Version</th><th class="text-center">Hardware</th></tr>
    </thead>
    <tbody>
      <tr>
        <td>{{ device.uid }}</td>
        <td class="text-center">{{ device.version }}</td>
        <td class="text-center">{{ device.hardware }}</td>
      </tr>
    </tbody>
  </table>
  </template>

  <div class="card-group">
    <div class="card card-body">
      <strong class="card-title">Device Identity</strong>

      <ul class="list-group list-group-flush">
        <li v-for="(value, key) in device.device_identity" class="list-group-item" :key="key">
          <span>{{ key }}</span><br/>
          <code>{{ value }}</code>
        </li>
      </ul>
    </div>

    <div class="card card-body">
      <strong class="card-title">Device Attributes</strong>

      <ul class="list-group list-group-flush">
        <li v-for="(value, key) in device.device_attributes" class="list-group-item" :key="key">
          <span>{{ key }}</span><br/>
          <code>{{ value }}</code>
        </li>
      </ul>
    </div>
  </div>
</div>
</template>

<script>
export default {
  name: 'DeviceDetail',

  props: ['uid', 'embedded'],

  data () {
    return {
      device: {}
    }
  },

  async created () {
    const uid = this.$route.params.uid || this.uid
    this.device = await this.getDevice(uid)
  },

  methods: {
    async getDevice (uid) {
      return this.$http.get('/api/devices/' + uid).then(res => {
        return res.data
      })
    }
  },

  filters: {
    pretty: function (value) {
      return JSON.stringify(value, null, 2)
        .replace(/\n/g, '<br/>')
        .replace(/ /g, '&nbsp;')
    }
  }
}
</script>
