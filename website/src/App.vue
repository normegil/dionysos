<template>
  <div class="app">
    <Navigation />
    <Content />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Navigation from "./layout/navigation/Navigation.vue";
import Content from "./layout/content/Content.vue";
@Component({
  components: { Content, Navigation }
})
export default class App extends Vue {}
</script>

<style lang="scss">
.app {
  color: $color-grey-dark;
  display: flex;
  height: 100vh;
  font-size: 1.6rem;
}
</style>
