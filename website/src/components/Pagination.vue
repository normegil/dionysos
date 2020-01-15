<template>
  <div class="pagination">
    <a class="pagination__button pagination__button--left"
      ><i class="las la-angle-double-left"/></a
    ><a class="pagination__button"><i class="las la-angle-left"/></a
    ><a class="pagination__button"><i class="las la-angle-right"/></a
    ><a class="pagination__button pagination__button--right"
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
  numberOfItems!: number;

  @Prop({ required: true })
  itemPerPage!: number;

  @Prop({ required: true })
  currentIndex!: number;
}
</script>

<style lang="scss">
.pagination {
  &__button {
    display: inline-block;
    padding: 0.5rem;
    border: 1px solid $color-grey-light-2;
    text-decoration: none;
    color: $color-grey-dark-3;
    transition: all 0.3s;

    &:hover {
      border: 1px solid $color-grey-light-1;
    }
  }
}
</style>
