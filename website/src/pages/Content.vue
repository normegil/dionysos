<template>
  <div class="content">
    <SearchField />
    <div class="content__title-container">
      <div class="content__title-background" />
      <div class="content__title-text">
        <h1>{{ $t("ui.menu.main.items") }}</h1>
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

  &__title {
    &-container {
      overflow: hidden;
      height: 9rem;
      position: relative;
    }

    &-background {
      background-image: url("../assets/images/items-title.jpg");
      background-size: cover;
      background-position: center;
      filter: blur(8px);
      height: 100%;
    }

    &-text {
      color: $color-black;
      position: absolute;
      top: 50%;
      left: 20rem;
    }
  }
}
</style>
