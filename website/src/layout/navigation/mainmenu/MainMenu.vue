<template>
  <div class="main-menu__container">
    <MainMenuItem v-for="link in links" :key="link.title" :link="link" />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import MainMenuItem from "./MainMenuItem.vue";
import LinkWithIcon from "../../../model/LinkWithIcon";
@Component({
  components: { MainMenuItem }
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
}
</script>

<style lang="scss">
.main-menu {
}
</style>
