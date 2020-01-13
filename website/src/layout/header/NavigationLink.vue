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
  text-decoration: none;
  color: $color-secondary;
  font-size: 1.6rem;
  display: inline-block;
  margin: 0 5rem;
}
</style>
