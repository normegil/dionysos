<template>
  <div class="main-menu__container">
    <MainMenuItem v-for="link in links" :key="link.title" :link="link" />
    <UnfoldableMenuItem
      title="Help"
      icon="las la-question"
      :sub-items="subItems"
    />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import MainMenuItem from "./MainMenuItem.vue";
import LinkWithIcon from "../../../model/LinkWithIcon";
import UnfoldableMenuItem from "./UnfoldableMenuItem.vue";
@Component({
  components: { UnfoldableMenuItem, MainMenuItem }
})
export default class Navigation extends Vue {
  links: LinkWithIcon[] = [
    {
      title: "Storages",
      href: "#",
      icon: "las la-warehouse"
    },
    {
      title: "Items",
      href: "#",
      icon: "las la-boxes"
    }
  ];

  subItems: LinkWithIcon[] = [
    {
      title: "Documentation",
      href: "#",
      icon: "las la-book"
    }
  ];
}
</script>

<style lang="scss">
.main-menu {
}
</style>
