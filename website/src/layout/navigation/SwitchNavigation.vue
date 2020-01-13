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
  font-size: 3rem;
  color: $color-grey-dark;

  &:hover {
    color: $color-grey-light-2;
  }
}
</style>
