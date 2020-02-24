<template>
  <div class="search-result__container-box">
    <h2 class="search-result__title">{{ result.type }}</h2>
    <ul class="search-result__list">
      <li
        class="search-result__list-item"
        v-for="resultItem in result.results"
        :key="resultItem.id"
      >
        {{ resultItem.name }}
      </li>
    </ul>
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";
import SearchResult, { PrintableResult } from "../model/SearchResult";
import Table from "./Table.vue";
@Component({
  components: { Table }
})
export default class SearchField extends Vue {
  @Prop({ required: true })
  result!: SearchResult<PrintableResult>;
}
</script>

<style lang="scss">
.search-result {
  &__container-box {
    margin-top: 2.2rem;
    padding: 1rem;
    border-radius: 2px;
    box-shadow: 0 0 5px $color-grey-light-2;
  }

  &__title {
    font-family: "Rokkitt", "Ubuntu Mono", sans-serif;
    font-weight: 400;
    font-size: 2.4rem;
    margin-bottom: 0.5rem;
  }

  &__list {
    list-style: none;
  }

  &__list-item {
    display: block;
    padding: 1rem 1rem;

    &:hover {
      background-color: $color-grey-light-2;
    }
  }
}
</style>
