<template>
  <div class="collection__container">
    <div class="collection__container-box flex-content">
      <Pagination
        :current-page="currentPage"
        :item-per-page="itemsPerPage"
        :number-of-pages="numberOfPages"
        @change-page="changePage"
      />
      <SpecificSearchField
        class="collection__filter"
        :placeholder="filterPlaceholder"
        :searched="filter"
        @search="filterCollection"
      />
      <Button
        icon="las la-plus"
        :title="$t('ui.button.add')"
        @click="$emit('create-item')"
      />
    </div>
    <div class="collection__container-box">
      <Table
        :headings="tableHeaders"
        :content="items"
        @edit="editItem"
        @remove="removeItem"
      />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import Button from "./Button.vue";
import Table from "./Table.vue";
import Rowable from "../model/Rowable";
import TableColumn from "../model/TableColumn";
import SpecificSearchField from "./SpecificSearchField.vue";
import { AxiosError } from "axios";
import Error from "../model/Error";
import Pagination from "./Pagination.vue";
@Component({
  components: { Pagination, SpecificSearchField, Table, Button }
})
export default class CollectionManager extends Vue {
  @Prop({ required: true })
  storeNamespace!: string;

  @Prop({ required: true })
  tableHeaders!: TableColumn[];

  @Prop({ default: "", required: false })
  filterPlaceholder!: string;

  @Prop({ default: 20, required: false })
  itemsPerPage!: number;

  get items(): Rowable[] {
    return this.$store.state[this.storeNamespace].items;
  }

  get filter(): string {
    return this.$store.state[this.storeNamespace].filter;
  }

  set filter(filter: string) {
    this.$store.commit(this.storeNamespace + "/setFilter", filter.trim());
  }

  get totalNumberOfItems(): number {
    return this.$store.state[this.storeNamespace].total;
  }

  get currentIndex(): number {
    return this.$store.state[this.storeNamespace].currentIndex;
  }

  set currentIndex(currentIndex: number) {
    let toSet = currentIndex;
    const i = this.lastPageFirstIndex;
    if (i < currentIndex) {
      toSet = i;
    }
    if (currentIndex < 0) {
      toSet = 0;
    }
    this.$store.commit(this.storeNamespace + "/setCurrentIndex", toSet);
  }

  get currentPage(): number {
    return Math.floor(this.currentIndex / this.itemsPerPage);
  }

  get lastPageFirstIndex(): number {
    const nbPage = this.numberOfPages;
    let lastPage = nbPage - 1; // pages start at 0
    if (nbPage === 0) {
      lastPage = 0;
    }
    return this.getPageFirstIndex(lastPage);
  }

  get isPageFullyContained(): boolean {
    return this.totalNumberOfItems % this.itemsPerPage === 0;
  }

  get numberOfPages(): number {
    let numberOfPages = Math.floor(this.totalNumberOfItems / this.itemsPerPage);
    if (!this.isPageFullyContained) {
      numberOfPages += 1;
    }
    console.log(this.totalNumberOfItems);
    return numberOfPages;
  }

  getPageFirstIndex(page: number): number {
    return this.itemsPerPage * page;
  }

  changePage(page: number): void {
    this.currentIndex = page * this.itemsPerPage;
    this.$store
      .dispatch(this.storeNamespace + "/load")
      .catch((err: AxiosError<Error>) => {
        this.error(err);
      });
  }

  editItem(id: string): void {
    this.$emit("edit-item", id);
  }

  removeItem(id: string): void {
    this.$emit("remove-item", id);
  }

  filterCollection(filter: string): void {
    this.filter = filter;
    this.currentIndex = 0;
    this.$store
      .dispatch(this.storeNamespace + "/load")
      .catch((err: AxiosError<Error>) => {
        this.error(err);
      });
  }

  mounted(): void {
    this.$store.dispatch(this.storeNamespace + "/load");
  }

  error(err: AxiosError<Error>): void {
    console.log(err);
  }
}
</script>

<style lang="scss">
.collection {
  &__container-box {
    margin-top: 2.2rem;
    padding: 1rem;
    border-radius: 2px;
    box-shadow: 0 0 5px $color-grey-light-2;
  }

  &__filter {
    width: 100%;
    margin: 0 1rem;
  }
}
.flex-content {
  display: flex;
}
</style>
