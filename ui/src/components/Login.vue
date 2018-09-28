<template>
<div class="container-fluid h-100 py-5" style="background-color:#0e293e">
  <div class="container">
  <div class="row">
    <div class="col-md-12">
      <p class="text-center mb-4"><img src="https://www.updatehub.io/imgs/updatehub-logo.png" width=200/></p>
      <div class="row">
        <div class="col-md-4 mx-auto">
          <div class="card rounded-4">
            <div class="card-body">
              <form class="form" role="form" autocomplete="off" novalidate="">
                <div class="form-group">
                  <label for="uname1">Username</label>
                  <input v-model="username" type="text" class="form-control form-control-lg rounded-0" name="uname1" required="">
                  <div class="invalid-feedback">Oops, you missed this one.</div>
                </div>
                <div class="form-group">
                  <label>Password</label>
                  <input v-model="password" type="password" class="form-control form-control-lg rounded-0" required="" autocomplete="new-password">
                  <div class="invalid-feedback">Enter your password too!</div>
                </div>
                <button v-on:click="login" class="btn btn-success btn-block btn-lg">Login</button>
              </form>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
</div>
</template>

<script>
export default {
  name: "Login",

  data() {
    return {
      username: "",
      password: ""
    };
  },

  methods: {
    async login(e) {
      let form = { username: this.username, password: this.password };
      this.$http
        .post("/login", form)
        .then(res => {
          this.$app.currentUser = res.data;
          this.$router.push(this.$route.query.redirect);
        })
        .catch(e => {
          console.log(e);
        });

      e.preventDefault();
    }
  }
};
</script>

<style>
</style>
