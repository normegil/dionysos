<template>
  <div class="content">
    <SearchField />
    <div class="content__main">
      <div class="content__title-container">
        <h1 class="content__title-text">{{ $t("ui.menu.main.items") }}</h1>
      </div>
    </div>
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

  &__main {
    margin-left: 3.7rem;
  }

  &__title {
    &-text {
      font-family: "Rokkitt", "Ubuntu Mono", sans-serif;
      font-weight: 400;
      margin-top: 2.2rem;
      font-size: 2.8rem;
    }
  }
}
</style>
