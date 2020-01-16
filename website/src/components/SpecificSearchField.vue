<template>
  <div class="specific-search-field">
    <!--suppress HtmlFormInputWithoutLabel -->
    <input
      class="specific-search-field__input"
      type="text"
      :placeholder="$t(placeholderKey)"
      v-model.trim="searched"
      @keyup.enter="search"
    />
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";

@Component
export default class SpecificSearchField extends Vue {
  @Prop({ required: true })
  placeholderKey!: string;

  searched = "";

  search(): void {
    console.log("Searching (specific): " + this.searched);
  }
}
</script>

<style lang="scss">
.specific-search-field {
  &__input {
    display: inline-block;
    color: $color-grey-dark-3;
    padding: 0.5rem;
    width: 100%;
    border: 1px solid $color-grey-light-2;
    border-radius: 5px;
    font-size: 1.6rem;
    transition: all 0.3s;

    &:hover,
    &:focus {
      border: 1px solid $color-grey-light-1;
    }
  }
}
</style>
