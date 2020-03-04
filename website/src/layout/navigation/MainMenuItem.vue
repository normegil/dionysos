<template>
  <div v-if="allowed">
    <router-link v-if="!isLink" :to="link.href" class="main-menu-item">
      <i :class="link.icon" class="main-menu-item__icon" />
      <span class="main-menu-item__text"> {{ link.title }}</span>
    </router-link>
    <a v-if="isLink" :href="link.href" class="main-menu-item">
      <i :class="link.icon" class="main-menu-item__icon" />
      <span class="main-menu-item__text"> {{ link.title }}</span>
    </a>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import MainLink from "../../model/MainLink";

@Component
export default class MainMenuItem extends Vue {
  @Prop({ default: {}, required: true })
  link!: MainLink;

  @Prop({ default: false, required: false })
  isLink!: boolean;

  get allowed(): boolean {
    return this.$store.getters["auth/hasAccess"](this.link.resource, "read");
  }
}
</script>

<style lang="scss">
.main-menu-item {
  display: block;
  font-size: 1.9rem;
  padding: 0.75rem 0;
  border-bottom: 1px solid $color-grey-dark-3;
  padding-left: 2.2rem;
  color: $color-grey-dark;
  text-decoration: none;
  transition: all 0.3s;

  &:hover {
    background-color: $color-grey-dark-3;
    color: $color-white;
  }

  &__text {
    display: inline-block;
    margin-left: 0.75rem;
  }
}
</style>
