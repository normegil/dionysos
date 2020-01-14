<template>
  <a :href="link.href" class="bottom-menu-link">{{ link.title }}</a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import Link from "../../../model/Link";

@Component
export default class BottomMenuLink extends Vue {
  @Prop({ default: {}, required: true })
  link!: Link;
}
</script>

<style lang="scss">
.bottom-menu-link {
  &:link,
  &:visited {
    font-size: 1.5rem;
    text-transform: uppercase;
    color: $color-grey-dark;
    text-decoration: none;
  }

  &:hover,
  &:active {
    color: $color-white;
  }
}
</style>
