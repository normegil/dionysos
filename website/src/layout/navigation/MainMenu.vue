<template>
  <div class="main-menu">
    <div class="top-menu__container">
      <MainMenuItem v-for="link in links" :key="link.title" :link="link" />
    </div>
    <div class="bottom-menu">
      <UnfoldableMenuItem
        :title="$t('ui.menu.main.more.title')"
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
import UnfoldableMenuItem from "../../components/UnfoldableMenuItem.vue";

@Component({
  components: { MainMenuItem, UnfoldableMenuItem }
})
export default class MainMenu extends Vue {
  get links(): LinkWithIcon[] {
    return [
      {
        title: this.$t("ui.menu.main.storages") as string,
        href: "/storages",
        icon: "las la-warehouse"
      },
      {
        title: this.$t("ui.menu.main.items") as string,
        href: "/items",
        icon: "las la-boxes"
      }
    ];
  }

  get moreMenuItems(): LinkWithIcon[] {
    return [
      {
        title: this.$t("ui.menu.main.more.sub.documentation") as string,
        href: "#",
        icon: "las la-book"
      },
      {
        title: "API",
        href: "api/",
        icon: "las la-server"
      },
      {
        title: this.$t("ui.menu.main.more.sub.bugs") as string,
        href: "#",
        icon: "las la-bug"
      },
      {
        title: "Github",
        href: "#",
        icon: "lab la-github"
      },
      {
        title: this.$t("ui.menu.main.more.sub.licenses") as string,
        href: "#",
        icon: "lab la-creative-commons-by"
      }
    ];
  }
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
