<template>
  <div>
    <div class="pb-0 mt-0 mb-4 border-bottom">
      <h3>
        Package Details
        <router-link to="/packages" class="btn btn-link float-sm-right">
          <i class="fas fa-trash"></i> Delete
        </router-link>
        <router-link to="/packages" class="btn btn-link float-sm-right">
          <i class="fas fa-caret-left"></i> Back to Packages
        </router-link>
      </h3>
    </div>

    <table class="table">
      <thead class="thead-light">
        <tr>
          <th>UID</th>
          <th>Version</th>
          <th>Supported Hardware</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{ pkg.uid }}</td>
          <td>{{ pkg.version }}</td>
          <td>{{ pkg.supported_hardware.join(', ') }}</td>
        </tr>
      </tbody>
    </table>

    <div class="card-group">
      <div class="card card-body">
        <strong class="card-title">Objects</strong>

        <table class="table table-hover">
          <thead>
            <tr>
              <th>Name</th>
              <th>Checksum</th>
              <th>Size</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(object, i) in pkg.metadata.objects[0]" :key="i" @click="selectedObject = i">
              <td>
                <span>{{ object.filename }}</span>
              </td>
              <td>1234...</td>
              <td>{{ object.size }}</td>
            </tr>
          </tbody>
        </table>
      </div>
      <div class="card" v-if="selectedObject >= 0">
        <div class="card-header">
          <ul class="nav nav-tabs card-header-tabs nav-fill">
            <li v-for="i in [0, 1]" :key="i" class="nav-item">
              <a
                class="nav-link"
                v-bind:class="currentInstallSet == i ? 'active' : ''"
                @click="currentInstallSet = i"
              >Installation Set #{{ i }}</a>
            </li>
          </ul>
        </div>
        <div class="card-body">
          <ul class="list-group list-group-flush">
            <li class="list-group-item">
              <span>Mode</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject].mode }}</code>
            </li>
            <li class="list-group-item">
              <span>Target</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject].target }}</code>
            </li>
            <li class="list-group-item">
              <span>Target Path</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject]['target-path'] }}</code>
            </li>
            <li class="list-group-item">
              <span>Target Type</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject]['target-type'] }}</code>
            </li>
            <li class="list-group-item">
              <span>Filesystem</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject].filesystem }}</code>
            </li>
            <li class="list-group-item">
              <span>Format Device</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject]['format?'] }}</code>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'PackageDetails',

  data () {
    return {
      pkg: { supported_hardware: [], metadata: { objects: [[]] } },
      selectedObject: -1,
      currentInstallSet: 0
    }
  },

  async created () {
    this.pkg = await this.getPackage()
    this.pkg.metadata = JSON.parse(window.atob(this.pkg.metadata))
  },

  methods: {
    async getPackage () {
      return this.$http
        .get('/api/packages/' + this.$route.params.uid)
        .then(res => {
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

<style scoped>
tr {
  cursor: pointer;
}

a.nav-link {
  cursor: pointer;
}
</style>
