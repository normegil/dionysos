<template>
  <div class="app">
    <Navigation class="app__navigation" />
    <router-view class="app__content" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Navigation from "./layout/navigation/Navigation.vue";
import Content from "./pages/PageItem.vue";
@Component({
  components: { Content, Navigation }
})
export default class App extends Vue {}
</script>

<style lang="scss">
.app {
  color: $color-grey-dark-3;
  display: flex;
  height: 100vh;
  font-size: 1.6rem;
  overflow: hidden;

  &__content {
    width: 100%;
  }

  &__navigation {
    min-width: 19rem;
  }
}
</style>
