<template>
<div class="container-fluid h-100 py-5">
  <div class="container">
  <div class="row">
    <div class="col-md-12">
      <p class="text-center mb-4"><img src="https://www.updatehub.io/imgs/updatehub-logo.png" width=200/></p>
      <div class="row">
        <div class="col-md-4 mx-auto">
          <div class="card rounded-4">
            <div class="card-body">
              <form class="form" role="form" @submit.prevent="login">
                <div class="form-group">
                  <label>Username</label>
                  <input v-model="username" type="text" class="form-control form-control-lg rounded-0" required v-focus/>
                </div>
                <div class="form-group">
                  <label>Password</label>
                  <input v-model="password" type="password" class="form-control form-control-lg rounded-0" required/>
                </div>
                <div class="alert alert-danger text-center" v-if="error">Invalid username or password</div>
                <button type="submit" class="btn btn-success btn-block btn-lg">Login</button>
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
  name: 'Login',

  data () {
    return {
      error: false,
      username: '',
      password: ''
    }
  },

  methods: {
    async login (e) {
      const form = { username: this.username, password: this.password }
      this.$http
        .post('/login', form)
        .then(res => {
          this.$app.currentUser = res.data
          this.error = false
          if (this.$route.query.redirect) {
            this.$router.push(this.$route.query.redirect)
          } else {
            this.$router.push('/')
          }
        })
        .catch(e => {
          this.error = true
        })

      e.preventDefault()
    }
  },

  directives: {
    focus: {
      inserted: function (el) {
        el.focus()
      }
    }
  }
}
</script>

<style scoped>
.container-fluid {
  background-color: #0e293e;
}
</style>
