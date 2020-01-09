<template>
  <ul>
    <li v-for="link in links" :key="link.title">
      <NavigationLink :link="link" />
    </li>
  </ul>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import Link from "../../model/Link";
import NavigationLink from "./NavigationLink.vue";
@Component({
  components: { NavigationLink }
})
export default class NavigationBar extends Vue {
  @Prop({ default: [], required: true })
  links!: Link[];
}
</script>

<style scoped></style>
