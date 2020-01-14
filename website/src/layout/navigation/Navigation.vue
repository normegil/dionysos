<template>
  <div class="navigation">
    <div class="navigation__logo-container">
      <Logo class="navigation__logo" />
      <SwitchNavigation />
    </div>
    <MainMenu />
    <BottomMenu class="navigation__bottom-menu" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Logo from "./Logo.vue";
import SwitchNavigation from "./SwitchNavigation.vue";
import MainMenu from "./mainmenu/MainMenu.vue";
import BottomMenu from "./bottommenu/BottomMenu.vue";

@Component({
  components: { BottomMenu, MainMenu, SwitchNavigation, Logo }
})
export default class Navigation extends Vue {}
</script>

<style lang="scss">
.navigation {
  background-image: linear-gradient(
      rgba($color-black, 0.85),
      rgba($color-black, 0.95)
    ),
    url("../../assets/images/navigation-background.jpg");
  background-size: cover;
  color: $color-grey-light-2;
  width: 60rem;
  position: relative;

  &__logo-container {
    background-color: $color-black;
    display: flex;
    align-items: baseline;
    padding: 0.8rem;
  }

  &__logo {
    margin-right: 5rem;
    margin-left: 1rem;
  }

  &__bottom-menu {
    position: absolute;
    bottom: 0;
  }
}
</style>
