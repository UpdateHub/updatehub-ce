<template>
  <div>
    <div class="pb-0 mt-0 mb-4 border-bottom">
      <h3>
        Package Details
        <a class="btn btn-link float-sm-right" @click="deletePackage">
          <i class="fas fa-trash"></i> Delete
        </a>
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
          <td>{{ [pkg.supported_hardware].join(', ') }}</td>
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
          <ul
            class="list-group list-group-flush"
            v-for="(item, index) in items"
            :key="index"
            >
            <li
              class="list-group-item"
              v-if="pkg.metadata.objects[currentInstallSet][selectedObject][item.field]"
              >
              <span>{{ item.title }}</span>
              <br>
              <code>{{ pkg.metadata.objects[currentInstallSet][selectedObject][item.field] }}</code>
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
      currentInstallSet: 0,
      items: [
        {
          title: 'Filename',
          field: 'filename'
        },
        {
          title: 'Compressed',
          field: 'compressed'
        },
        {
          title: 'Required Uncompressed Size',
          field: 'required-uncompressed-size'
        },
        {
          title: 'Chunk Size',
          field: 'chunk-size'
        },
        {
          title: 'Skip',
          field: 'skip'
        },
        {
          title: 'Seek',
          field: 'seek'
        },
        {
          title: 'Count',
          field: 'count'
        },
        {
          title: 'Truncate',
          field: 'truncate'
        },
        {
          title: 'Filesystem',
          field: 'filesystem'
        },
        {
          title: 'Target',
          field: 'target'
        },
        {
          title: 'Target Path',
          field: 'target-path'
        },
        {
          title: 'Format',
          field: 'format?'
        },
        {
          title: 'Format Options',
          field: 'format-options'
        },
        {
          title: 'Mount Options',
          field: 'mount-options'
        },
        {
          title: '1K Padding',
          field: '1k_padding'
        },
        {
          title: 'Search Exponent',
          field: 'search_exponent'
        },
        {
          title: 'Chip 0 Device Path',
          field: 'chip_0_device_path'
        },
        {
          title: 'Chip 1 Device Path',
          field: 'chip_1_device_path'
        }
      ]
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
    },

    deletePackage () {
      this.$http
        .delete('/api/packages/' + this.$route.params.uid + '/delete')
        .then(res => {
          this.$router.push('/packages')
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
