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
      packages: []
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
      let form = new FormData();
      form.append("file", this.$refs.file.files[0]);

      this.$http
        .post("/api/packages", form, {
          headers: { "Content-Type": "multipart/form-data" }
        })
        .then(function() {
          console.log("SUCCESS!!");
        })
        .catch(function() {
          console.log("FAILURE!!");
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