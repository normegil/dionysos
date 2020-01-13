<template>
  <a class="navigation-link" :href="link.href">{{ link.title }}</a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import Link from "../../model/Link";

@Component
export default class NavigationLink extends Vue {
  @Prop({ default: {}, required: true })
  link!: Link;
}
</script>

<style lang="scss">
.navigation-link {
  &:link,
  &:visited {
    text-decoration: none;
    text-transform: capitalize;
    font-size: 1.6rem;
    display: inline-block;
    padding: 1rem;
    border-radius: 5px;
    margin: 0 1rem;
    transition: all 0.2s;
    color: $color-grey-dark;
  }

  &:hover,
  a:active {
    background-color: $color-primary-light;
    color: $color-grey-light-2;
  }
}
</style>
