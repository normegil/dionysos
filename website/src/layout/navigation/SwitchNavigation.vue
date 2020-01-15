<template>
  <div>
    <i class="las la-angle-left switch-navigation__icon" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

@Component
export default class SwitchNavigation extends Vue {}
</script>

<style lang="scss">
.switch-navigation__icon {
  font-size: 2.2rem;
  color: $color-grey-dark;
  transition: all 0.3s;

  &:hover {
    color: $color-white;
  }
}
</style>
