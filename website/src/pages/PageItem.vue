<template>
  <div class="content">
    <SearchField />
    <div class="content__main">
      <div class="content__title-container">
        <h1 class="content__title-text">{{ $t("ui.menu.main.items") }}</h1>
      </div>
      <div class="content__container flex-content">
        <Pagination
          :current-page="currentPage"
          :item-per-page="itemsPerPage"
          :number-of-pages="numberOfPages"
          @change-page="changePage"
        />
        <SpecificSearchField
          class="content__search-field"
          placeholder-key="ui.components.specific-search-field.items.placeholder"
        />
        <Button icon="las la-plus" />
      </div>
      <div class="content__container">
        <Table :headings="headings" :content="items" />
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import SearchField from "../components/SearchField.vue";
import Table from "../components/Table.vue";
import Item from "../model/Item";
import Button from "../components/Button.vue";
import Pagination from "../components/Pagination.vue";
import SpecificSearchField from "../components/SpecificSearchField.vue";
@Component({
  components: { SpecificSearchField, Pagination, Button, Table, SearchField }
})
export default class PageItem extends Vue {
  get items(): Item[] {
    return this.$store.state.items.items;
  }

  get numberOfPages(): number {
    return this.$store.getters["items/numberOfPages"];
  }

  get currentPage(): number {
    return this.$store.getters["items/currentPage"];
  }

  get itemsPerPage(): number {
    return this.$store.state.items.itemsPerPage;
  }

  get headings(): string[] {
    return [this.$t("ui.components.table.items.heading.name") as string];
  }

  mounted(): void {
    this.$store.dispatch("items/load");
  }

  changePage(pageNb: number): void {
    this.$store.dispatch("items/changePage", pageNb);
  }
}
</script>

<style lang="scss">
.content {
  background-color: $color-white;

  overflow: auto;

  &__main {
    margin: 0 3.7rem;
  }

  &__title {
    &-text {
      font-family: "Rokkitt", "Ubuntu Mono", sans-serif;
      font-weight: 400;
      margin-left: 2rem;
      margin-top: 2.2rem;
      font-size: 2.8rem;
    }
  }

  &__container {
    margin-top: 2.2rem;
    padding: 1rem;
    border-radius: 2px;
    box-shadow: 0 0 5px $color-grey-light-2;
  }
  &__search-field {
    width: 100%;
    margin: 0 1rem;
  }
}
.flex-content {
  display: flex;
}
</style>
