<template>
  <table class="table">
    <tr class="table__row-heading">
      <th class="table__heading" v-for="heading in headings" :key="heading">
        {{ heading }}
      </th>
    </tr>
    <tr class="table__row" v-for="row in content" :key="row.identifier()">
      <td
        class="table__cell"
        v-for="(cell, index) in row.toRow()"
        :key="row.identifier() + index"
      >
        {{ cell }}
      </td>
    </tr>
  </table>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";
import Rowable from "../model/Rowable";

@Component
export default class Table extends Vue {
  @Prop({ default: [], required: false })
  headings!: string[];

  @Prop({ default: [], required: true })
  content!: Rowable[];
}
</script>

<style lang="scss">
.table {
  width: 100%;
  border-collapse: collapse;

  &__heading {
    text-align: left;
    border-bottom: 2px solid $color-grey-dark;
    padding: 1rem;

    &:not(:last-child) {
      border-right: 2px solid $color-grey-dark;
    }
  }

  &__cell {
    padding: 1rem;
  }

  &__row {
    transition: all 0.3s;
    &:hover {
      background-color: $color-grey-light-2;
    }
  }

  &__row:not(:last-child) &__cell {
    border-bottom: 1px solid $color-grey-light-2;
  }
}
</style>
