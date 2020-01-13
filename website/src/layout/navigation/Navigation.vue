<template>
  <div class="navigation">
    <div class="navigation__main">
      <Logo class="navigation__logo" />
      <SwitchNavigation />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Logo from "./Logo.vue";
import SwitchNavigation from "./SwitchNavigation.vue";

@Component({
  components: { SwitchNavigation, Logo }
})
export default class Navigation extends Vue {}
</script>

<style lang="scss">
.navigation {
  background-color: $color-grey-dark-3;
  color: $color-grey-light-2;
  width: 60rem;

  &__main {
    background-color: $color-black;
    display: flex;
    align-items: baseline;
    padding: 0.8rem;
  }

  &__logo {
    margin-right: 2rem;
  }
}
</style>
