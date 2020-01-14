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
  background-image: linear-gradient(
      rgba($color-black, 0.85),
      rgba($color-black, 0.95)
    ),
    url("../../assets/images/navigation-background.jpg");
  background-size: cover;
  color: $color-grey-light-2;
  width: 60rem;

  &__main {
    background-color: $color-black;
    display: flex;
    align-items: baseline;
    padding: 0.8rem;
  }

  &__logo {
    margin-right: 5rem;
    margin-left: 2rem;
  }
}
</style>
