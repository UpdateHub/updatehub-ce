<template>
<div>
  <div class="pb-0 mt-0 mb-4 border-bottom">
    <h3>
      Package List
      <label class="btn btn-link float-sm-right">
        <i class="fas fa-upload"></i> Upload Package<input v-on:change="uploadPackage" type="file" ref="file" hidden/>
      </label>
    </h3>
  </div>

  <div class="progress mb-4" v-if="uploadProgress > 0">
    <div class="progress-bar" v-bind:style="{width: uploadProgress + '%'}"></div>
  </div>

  <div class="alert alert-danger" role="alert" v-if="lastError">
    <strong>Failed to upload package:</strong> {{ lastError }}
  </div>

  <table class="table table-bordered table-hover">
    <thead class="thead-light">
      <tr>
        <th>UID</th><th>Version</th><th>Supported Hardware</th>
      </tr>
    </thead>
    <tbody>
        <tr v-for="pkg in packages" :key="pkg.id" @click="$router.push(`/packages/${pkg.uid}`)">
          <td>{{ pkg.uid }}</td>
          <td>{{ pkg.version }}</td>
          <td>{{ pkg.supported_hardware.join(', ') }}</td>
        </tr>
    </tbody>
  </table>
</div>

</template>

<script>
export default {
  name: "PackageList",

  data() {
    return {
      packages: [],
      uploadProgress: 0,
      lastError: null,
    };
  },

  async created() {
    this.packages = await this.getPackages();
  },

  methods: {
    async getPackages() {
      return await this.$http.get("/api/packages").then(res => {
        return res.data;
      });
    },

    uploadPackage(e) {
      this.lastError = null;

      let form = new FormData();
      form.append("file", this.$refs.file.files[0]);

      this.$http
        .post("/api/packages", form, {
          headers: { "Content-Type": "multipart/form-data" },
          onUploadProgress: function(e) {
            this.uploadProgress = parseInt(Math.round((e.loaded * 100) / e.total));
          }.bind(this)
        })
        .then(res => {
          this.$router.push("/packages/" + res.data.uid);
        })
        .catch(e => {
          this.lastError = e.response.data.message;
          this.uploadProgress = 0;
        });
    }
  }
};
</script>

<style scoped>
tr {
  cursor: pointer;
}
</style>