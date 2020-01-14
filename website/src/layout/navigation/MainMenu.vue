<template>
  <div class="main-menu">
    <div class="top-menu__container">
      <MainMenuItem v-for="link in links" :key="link.title" :link="link" />
    </div>
    <div class="bottom-menu">
      <UnfoldableMenuItem
        title="More"
        icon="las la-ellipsis-h"
        :sub-items="moreMenuItems"
        direction="up"
      />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import MainMenuItem from "./MainMenuItem.vue";
import LinkWithIcon from "../../model/LinkWithIcon";
import UnfoldableMenuItem from "../components/UnfoldableMenuItem.vue";

@Component({
  components: { MainMenuItem, UnfoldableMenuItem }
})
export default class MainMenu extends Vue {
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

  moreMenuItems: LinkWithIcon[] = [
    {
      title: "Documentation",
      href: "#",
      icon: "las la-book"
    },
    {
      title: "API",
      href: "#",
      icon: "las la-server"
    },
    {
      title: "Bugs & suggestions",
      href: "#",
      icon: "las la-bug"
    },
    {
      title: "Github",
      href: "#",
      icon: "lab la-github"
    },
    {
      title: "Licenses",
      href: "#",
      icon: "lab la-creative-commons-by"
    }
  ];
}
</script>

<style lang="scss">
.main-menu {
  height: 100%;
}
.bottom-menu {
  position: absolute;
  bottom: 0;
  width: 100%;
}
</style>
