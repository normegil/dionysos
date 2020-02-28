<template>
  <div>
    <div class="user-display" @click.stop="mainAction">
      <i :class="icon" class="user-display__icon" />
      <span class="user-display__text">{{ textDisplayed }}</span>
    </div>
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";

@Component
export default class UserDisplay extends Vue {
  get icon(): string {
    if (this.isAuthenticated) {
      return "las la-sign-out-alt";
    }
    return "las la-sign-in-alt";
  }

  get textDisplayed(): string {
    if (this.isAuthenticated) {
      return this.$store.state.auth.authentifiedUser.name;
    }
    return "Sign-in";
  }

  get isAuthenticated(): boolean {
    return this.$store.getters["auth/isAuthenticated"];
  }

  mainAction(): void {
    if (this.isAuthenticated) {
      this.$store.dispatch("auth/signOut");
    } else {
      this.$store.commit("auth/setShowLoginModal", true);
    }
  }
}
</script>

<style lang="scss">
.user-display {
  font-size: 1.9rem;
  padding: 0.75rem 0;
  padding-left: 2.2rem;
  color: $color-grey-light-2;
  background-color: $color-black;
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s;

  &:hover {
    color: $color-primary;
  }

  &__text {
    display: inline-block;
    margin-left: 0.75rem;
  }
}
</style>
