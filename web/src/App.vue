<template>
  <AdbForm/>
</template>

<script>
import AdbForm from './components/AdbForm.vue'

export default {
  name: 'App',
  components: {
    AdbForm
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
  width: 80%;
  margin: 40px auto 0 auto;
}
</style>
