<template>
  <div class="pagination">
    <a class="pagination__button" @click.stop="changePage(0)"
      ><i class="las la-angle-double-left"/></a
    ><a class="pagination__button" @click.stop="changePage(currentPage - 1)"
      ><i class="las la-angle-left"/></a
    ><a
      v-for="page in pages"
      :key="page"
      class="pagination__button pagination__button--number"
      :class="{ 'pagination__button--active': page === currentPage }"
      @click.stop="changePage(page)"
      >{{ page }}</a
    ><a class="pagination__button" @click.stop="changePage(currentPage + 1)"
      ><i class="las la-angle-right"/></a
    ><a class="pagination__button" @click.stop="changePage(lastPage)"
      ><i class="las la-angle-double-right"
    /></a>
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";

@Component
export default class Pagination extends Vue {
  @Prop({ required: true })
  numberOfPages!: number;

  @Prop({ required: true })
  currentPage!: number;

  @Prop({ default: 5, required: false })
  numberDisplayedPage = 5;

  get pages(): number[] {
    const limits = this.pageLimits();
    const numbers: number[] = [];
    for (let i = limits.start; i <= limits.end; i++) {
      numbers.push(i);
    }
    return numbers;
  }

  get lastPage(): number {
    return this.numberOfPages - 1;
  }

  pageLimits(): { start: number; end: number } {
    let startIndex = this.currentPage - 2;
    let endIndex = this.currentPage + 2;

    if (this.lastPage < this.numberDisplayedPage) {
      return {
        start: 0,
        end: this.lastPage
      };
    }

    if (startIndex < 0) {
      const underflow = -startIndex;
      endIndex += underflow;
      startIndex = 0;
    }

    if (endIndex > this.lastPage) {
      const overflow = endIndex - this.lastPage;
      endIndex = this.lastPage;
      startIndex -= overflow;
      if (startIndex < 0) {
        startIndex = 0;
      }
    }

    return {
      start: startIndex,
      end: endIndex
    };
  }

  changePage(pageNb: number): void {
    this.$emit("change-page", pageNb);
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
      padding: 0.5rem 1rem;
      color: $color-grey-dark;
    }

    &--active {
      font-weight: bold;
      color: $color-grey-dark-3;
    }

    &:hover {
      color: $color-grey-dark-3;
      background-color: $color-grey-light-2;
    }
  }
}
</style>
