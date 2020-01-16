<template>
  <div class="pagination">
    <a class="pagination__button"><i class="las la-angle-double-left"/></a
    ><a class="pagination__button"><i class="las la-angle-left"/></a
    ><a
      v-for="page in pages"
      :key="page"
      class="pagination__button pagination__button--number"
      >{{ page }}</a
    ><a class="pagination__button"><i class="las la-angle-right"/></a
    ><a class="pagination__button"><i class="las la-angle-double-right"/></a>
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";

@Component
export default class Pagination extends Vue {
  @Prop({ required: true })
  numberOfItems!: number;

  @Prop({ required: true })
  itemPerPage!: number;

  @Prop({ required: true })
  currentIndex!: number;

  get pages(): number[] {
    const numbers: number[] = [];
    for (let i = 2; i !== 0 && this.currentPage - i >= 0; i--) {
      numbers.push(this.currentPage - i);
    }
    numbers.push(this.currentPage);
    for (
      let i = 1;
      numbers.length < 5 && this.currentPage + i <= this.totalPages;
      i++
    ) {
      numbers.push(this.currentPage + i);
    }
    for (let i = 3; numbers.length < 5 && this.currentPage - i >= 0; i++) {
      numbers.push(this.currentPage - i);
    }
    return numbers.sort();
  }

  get currentPage(): number {
    let currentPage = Math.floor(this.currentIndex / this.itemPerPage);
    const overflow = this.currentIndex % this.itemPerPage;
    if (overflow !== 0) {
      currentPage += 1;
    }
    console.log("Current: " + currentPage);
    return currentPage;
  }

  get totalPages(): number {
    let numberOfPages = Math.floor(this.numberOfItems / this.itemPerPage);
    const divisionOverflow = this.numberOfItems % this.itemPerPage;
    if (divisionOverflow !== 0) {
      numberOfPages += 1;
    }
    console.log("Total: " + numberOfPages);
    return numberOfPages;
  }
}
</script>

<style lang="scss">
.pagination {
  white-space: nowrap;
  &__button {
    display: inline-block;
    padding: 0.5rem;
    border: 1px solid $color-grey-light-2;
    text-decoration: none;
    color: $color-grey-dark-3;
    transition: all 0.3s;

    &:not(:last-child) {
      margin-right: -1px;
    }

    &:first-child {
      border-top-left-radius: 5px;
      border-bottom-left-radius: 5px;
    }

    &:last-child {
      border-top-right-radius: 5px;
      border-bottom-right-radius: 5px;
    }

    &--number {
      padding: 0.56rem 1rem 0.5rem;
    }

    &:hover {
      border: 1px solid $color-grey-dark;
    }
  }
}
</style>
