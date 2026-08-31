<template>
  <div class="app-layout">
    <AdbForm />
  </div>
</template>

<script>
import AdbForm from './components/AdbForm.vue';

export default {
  name: 'App',
  components: {
    AdbForm
  }
};
</script>

<style>
body {
  margin: 0;
  padding: 0;
  background-color: #f0f2f5;
}

#app {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
  width: 92%;
  max-width: 1200px;
  min-width: 320px;
  margin: 24px auto;
  box-sizing: border-box;
}

@media (max-width: 768px) {
  #app {
    width: 96%;
    margin: 12px auto;
  }
}
</style>
