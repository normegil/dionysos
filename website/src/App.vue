<template>
  <Footer />
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Footer from "./layout/footer/Footer.vue";
@Component({
  components: { Footer }
})
export default class App extends Vue {}
</script>

<style lang="scss">
@import "src/assets/scss/colors";
</style>
