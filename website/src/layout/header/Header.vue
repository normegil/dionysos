<template>
  <header class="header">
    <Logo class="header__logo" />
    <NavigationBar :links="links" />
  </header>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import NavigationBar from "./NavigationBar.vue";
import Link from "../../model/Link";
import Logo from "./Logo.vue";

@Component({
  components: { NavigationBar, Logo }
})
export default class Header extends Vue {
  links: Link[] = [
    {
      title: "Home",
      href: "/"
    },
    {
      title: "API Documentation",
      href: "/api/"
    }
  ];
}
</script>

<style lang="scss">
.header {
  display: flex;
  align-items: center;

  background-color: $color-primary-darker;
  border-bottom: 1px solid $color-grey-light-1;
  padding: 0.4rem 0;

  &__logo {
    margin: 0 1rem;
  }
}
</style>
