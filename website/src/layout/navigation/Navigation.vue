<template>
  <div class="navigation">
    <Logo />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Logo from "./Logo.vue";

@Component({
  components: { Logo }
})
export default class Navigation extends Vue {}
</script>

<style lang="scss">
.navigation {
  background-color: $color-black;
  color: $color-grey-light-2;
  padding: 2rem;
  width: 60rem;
}
</style>
