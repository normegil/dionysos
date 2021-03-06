<template>
  <table class="table">
    <tr class="table__row-heading">
      <th
        class="table__heading"
        :style="heading.style"
        v-for="heading in headings"
        :key="heading.name"
      >
        {{ heading.name }}
      </th>
      <th class="table__heading table__heading-action">
        {{ $t("ui.components.table.items.heading.actions") }}
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
      <td class="table__cell table__cell-action">
        <Button
          v-if="!readonly"
          icon="las la-pen"
          @click="edit(row.identifier())"
        />
        <Button
          v-if="!readonly"
          icon="las la-trash"
          @click="remove(row.identifier())"
        />
      </td>
    </tr>
  </table>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";
import Rowable from "../model/Rowable";
import TableColumn from "../model/TableColumn";
import Button from "./Button.vue";
@Component({
  components: { Button }
})
export default class Table extends Vue {
  @Prop({ default: [], required: false })
  headings!: TableColumn[];

  @Prop({ default: [], required: true })
  content!: Rowable[];

  @Prop({ default: false, required: false })
  readonly!: boolean;

  edit(id: string): void {
    this.$emit("edit", id);
  }

  remove(id: string): void {
    this.$emit("remove", id);
  }
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

    &-action {
      width: 12rem;
    }
  }

  &__cell {
    padding: 1rem;
  }

  &__row {
    transition: all 0.3s;
    &:hover {
      background-color: $color-grey-light-3;
    }
  }

  &__row:not(:last-child) &__cell {
    border-bottom: 1px solid $color-grey-light-2;
  }
}
</style>
