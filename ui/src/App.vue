<template>
<div>
  <template v-if="$app.currentUser">
    <nav class="navbar navbar-expand-lg">
      <a class="navbar-brand">
        <img src="https://www.updatehub.io/imgs/updatehub-logo.png" width="200" />
      </a>
      <div class="collapse navbar-collapse">
        <ul class="navbar-nav mr-auto"></ul>
      </div>
    </nav>

    <ul class="nav nav-tabs nav-fill pt-5">
      <li class="nav-item">
        <router-link to="/overview" class="nav-link" active-class="active">
        <i class="fas fa-home"></i> Overview
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/devices" class="nav-link" active-class="active">
        <i class="fas fa-microchip"></i> Devices
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/packages" class="nav-link" active-class="active">
        <i class="fas fa-file-archive"></i> Packages
        </router-link>
      </li>
      <li class="nav-item">
        <router-link to="/rollouts" class="nav-link" active-class="active">
        <i class="fas fa-external-link-alt"></i> Rollouts
        </router-link>
      </li>
    </ul>

    <main class="pt-4">
      <div class="container">
        <router-view></router-view>
      </div>
    </main>
  </template>
  <template v-if="!$app.currentUser">
    <router-view></router-view>
  </template>
</div>
</template>

<script>
export default {
  name: "app"
};
</script>

<style scoped>
ul {
  background-color: #f7f7f7;
}

nav {
  padding: 10px;
  background-color: #0e293e;
}

main {
  background-color: #fff;
}
</style>
