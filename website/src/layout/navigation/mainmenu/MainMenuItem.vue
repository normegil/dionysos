<template>
  <a :href="link.href" class="main-menu-item">
    <i :class="link.icon" class="main-menu-item__icon" />
    <span class="main-menu-item__text"> {{ link.title }}</span>
  </a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import LinkWithIcon from "../../../model/LinkWithIcon";

@Component
export default class MainMenuItem extends Vue {
  @Prop({ default: {}, required: true })
  link!: LinkWithIcon;
}
</script>

<style lang="scss">
.main-menu-item {
  display: block;
  font-size: 2.5rem;
  padding: 1rem 0;
  border-bottom: 1px solid $color-grey-dark-3;
  padding-left: 3rem;
  color: $color-grey-dark;
  text-decoration: none;
  transition: all 0.3s;

  &:hover {
    background-color: $color-grey-dark-3;
    color: $color-white;
  }

  &__text {
    display: inline-block;
    margin-left: 1rem;
  }
}
</style>
