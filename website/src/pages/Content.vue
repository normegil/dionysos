<template>
  <div class="content">
    <SearchField />
    <div class="content__main">
      <div class="content__title-container">
        <h1 class="content__title-text">{{ $t("ui.menu.main.items") }}</h1>
      </div>
      <div class="content__container flex-content">
        <Pagination current-index="5" item-per-page="2" number-of-items="10" />
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
@Component({
  components: { Pagination, Button, Table, SearchField }
})
export default class Content extends Vue {
  items: Item[] = [
    new Item("0", "Tomates"),
    new Item("1", "Oignons"),
    new Item("2", "Carotte"),
    new Item("3", "Spaggethi")
  ];

  get headings(): string[] {
    return [this.$t("ui.components.table.items.heading.name") as string];
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
}
.flex-content {
  display: flex;
}
</style>
