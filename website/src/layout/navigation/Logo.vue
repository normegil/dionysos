<template>
  <a href="/" class="logo">Dionysos</a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";

@Component
export default class Logo extends Vue {}
</script>

<style lang="scss">
.logo {
  &:link,
  &:visited {
    text-decoration: none;
    font-family: "Rokkitt", "Ubuntu Mono", sans-serif;
    font-size: 3rem;
    font-weight: 300;
    color: $color-white;
    transition: all 0.3s;
  }

  &:hover,
  &:active {
    color: $color-primary;
  }
}
</style>
