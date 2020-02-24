<template>
  <div class="content">
    <SearchField />
    <div class="content__main">
      <div class="content__title-container">
        <h1 class="content__title-text">{{ $t("ui.menu.main.search") }}</h1>
      </div>
      <SearchResultView
        v-for="result in results"
        :key="result.type"
        :result="result"
      />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import SearchField from "../components/SearchField.vue";
import SearchResultView from "../components/SearchResultView.vue";
import SearchResult, { PrintableResult } from "../model/SearchResult";

@Component({
  components: {
    SearchResultView,
    SearchField
  }
})
export default class PageSearch extends Vue {
  get results(): SearchResult<PrintableResult> {
    return this.$store.state.searchResults;
  }
}
</script>

<style lang="scss">
.content {
  background-color: $color-white;

  overflow: auto;

  &__title {
    &-text {
      font-family: "Rokkitt", "Ubuntu Mono", sans-serif;
      font-weight: 400;
      margin-left: 2rem;
      margin-top: 2.2rem;
      font-size: 2.8rem;
    }
  }
}
</style>
