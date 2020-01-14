<template>
  <div class="bottom-menu">
    <BottomMenuLink
      class="bottom-menu__item"
      v-for="link in links"
      :key="link.title"
      :link="link"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import Link from "../../../model/Link";
import BottomMenuLink from "./BottomMenuLink.vue";
@Component({
  components: { BottomMenuLink }
})
export default class BottomMenu extends Vue {
  @Prop({ default: {}, required: true })
  links: Link[] = [
    {
      title: "Documentation",
      href: "#"
    },
    {
      title: "Bugs",
      href: "#"
    },
    {
      title: "Licenses",
      href: "#"
    },
    {
      title: "API",
      href: "/api/"
    },
    {
      title: "Github",
      href: "https://github.com"
    }
  ];
}
</script>

<style lang="scss">
.bottom-menu {
  padding: 1rem 2rem;
  text-align: center;
  border-top: 1px solid $color-grey-dark-3;

  &__item {
    display: inline-block;

    &:not(:last-child) {
      padding-right: 1rem;
    }
  }
}
</style>
