<template>
  <a href="/" class="logo__container">
    <img :src="logoSrc" alt="logo" class="logo__image" />
  </a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import logo from "../../assets/images/logo.svg";

@Component
export default class Logo extends Vue {
  logoSrc = logo;
}
</script>

<style lang="scss">
.logo {
  &__container {
  }
  &__image {
    width: 4rem;
    padding: 0.2rem;
    border-radius: 2px;
    background-color: $white;
  }
}
</style>
