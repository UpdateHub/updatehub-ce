<template>
  <div>
    <div class="pb-0 mt-0 mb-4 border-bottom">
      <h3>
        Create Rollout
        <router-link to="/rollouts" class="btn btn-link float-sm-right">
          <i class="fas fa-caret-left"></i> Back to Rollouts
        </router-link>
      </h3>
    </div>

    <div class="card card-body">
      <form>
        <div class="form-group">
          <label>Package Version</label>
          <select class="form-control" @change="selectAllDevices" v-model="selectedPackage">
            <option disabled selected>Choose Version</option>
            <option v-for="pkg in packages" v-bind:value="pkg" :key="pkg.uid">{{ pkg.version }}</option>
          </select>
          <small>This is the version you will have when this rollout be completed.</small>
        </div>
      </form>

      <template v-if="selectedPackage">
        <small class="text-right">
          <strong>{{ selectedDevices.length }}</strong> of
          <strong>{{ devices.length }}</strong> devices will be updated to
          <strong>{{ selectedPackage.version }}</strong>
        </small>
        <table class="table table-hover">
          <thead class="thead-light">
            <tr>
              <th></th>
              <th>UID</th>
              <th>Version</th>
              <th>Hardware</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <template v-for="device in compatibleDevices()">
              <tr :key="device.uid" @click="toggleOpened(device.uid)">
                <td>
                  <input
                    type="checkbox"
                    value
                    checked
                    @click.stop="toggleDeviceSelection(device.uid)"
                  >
                </td>
                <td>
                  <span>{{ device.uid }}</span>
                </td>
                <td>
                  <span>{{ device.version }}</span>
                </td>
                <td>
                  <span>{{ device.hardware }}</span>
                </td>
                <td>
                  <i
                    class="fas"
                    v-bind:class="{ 'fa-caret-up': opened == device.uid, 'fa-caret-down': opened != device.uid }"
                  ></i>
                </td>
              </tr>
              <tr :key="device.uid + 'details'" v-if="opened == device.uid">
                <td colspan="5">
                  <DeviceDetails :uid="device.uid" embedded="true"></DeviceDetails>
                </td>
              </tr>
            </template>
          </tbody>
        </table>

        <div>
          <button
            class="btn btn-primary float-sm-right ml-1"
            @click="create()"
            :disabled="selectedDevices.length == 0"
          >Create</button>
          <button class="btn btn-outline-primary float-sm-left" @click="refresh">Refresh</button>
        </div>
      </template>
    </div>
  </div>
</template>

<script>
import DeviceDetails from './DeviceDetails'

export default {
  name: 'RolloutNew',

  components: { DeviceDetails },

  data () {
    return {
      selectedPackage: null,
      packages: [],
      devices: [],
      opened: '',
      selectedDevices: []
    }
  },

  async created () {
    this.refresh()
  },

  methods: {
    async refresh () {
      this.devices = await this.getDevices()
      this.packages = await this.getPackages()
    },

    async getDevices () {
      return this.$http.get('/api/devices').then(res => {
        return res.data
      })
    },

    async getPackages () {
      return this.$http.get('/api/packages').then(res => {
        return res.data
      })
    },

    toggleOpened (uid) {
      if (this.opened === uid) {
        this.opened = ''
      } else {
        this.opened = ''
        this.$nextTick(() => {
          this.opened = uid
        })
      }
    },

    toggleDeviceSelection (uid) {
      let index = this.selectedDevices.findIndex((device, i) => {
        return device.uid === uid
      })

      if (index >= 0) {
        this.selectedDevices.splice(index, 1)
      } else {
        this.selectedDevices.push(
          this.compatibleDevices().find(d => {
            return d.uid === uid
          })
        )
      }
    },

    selectAllDevices () {
      this.selectedDevices = this.compatibleDevices()
    },

    compatibleDevices () {
      return this.devices.filter(device => {
        const differentVersion =
          device.version !== this.selectedPackage.version
        const supportsAnyHardware =
          this.selectedPackage.supported_hardware === 'any'
        const supportsDeviceHardware = this.selectedPackage.supported_hardware.includes(
          device.hardware
        )

        return (
          differentVersion && (supportsAnyHardware || supportsDeviceHardware)
        )
      })
    },

    create () {
      this.$http
        .post('/api/rollouts', {
          package: this.selectedPackage.uid,
          devices: this.selectedDevices.map(d => {
            return d.uid
          })
        })
        .then(res => {
          let rollout = res.data
          this.$router.push('/rollouts/' + rollout.id)
        })
    }
  }
}
</script>
