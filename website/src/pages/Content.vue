<template>
  <div class="content">
    <SearchField />
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import SearchField from "../components/SearchField.vue";
@Component({
  components: { SearchField }
})
export default class Content extends Vue {}
</script>

<style lang="scss">
.content {
  background-color: $color-white;

  overflow: auto;
}
</style>
