<template>
  <footer>
    <NavigationBar :links="links" />
  </footer>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import NavigationBar from "./NavigationBar.vue";
import Link from "../../model/Link";
@Component({
  components: { NavigationBar }
})
export default class Footer extends Vue {
  links: Link[] = [
    {
      title: "API Documentation",
      href: "/api/"
    }
  ];
}
</script>

<style></style>
