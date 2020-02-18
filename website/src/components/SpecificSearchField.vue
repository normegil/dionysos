<template>
  <div class="specific-search-field specific-search-field__container">
    <!--suppress HtmlFormInputWithoutLabel -->
    <input
      class="specific-search-field__input"
      type="text"
      :placeholder="$t(placeholderKey)"
      v-model.trim="searched"
      @keyup.enter="search"
    />
    <a class="specific-search-field__button" @click.stop="search">{{
      $t("ui.button.search")
    }}</a>
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
  &__container {
    display: flex;
    flex-direction: row;
    border: 1px solid $color-grey-light-2;
    border-radius: 5px;
    transition: all 0.3s;

    &:hover,
    &:focus {
      margin: 2px;
      border: 1px solid $color-grey-light-1;
    }
  }

  &__input {
    flex-grow: 2;
    border: none;

    color: $color-grey-dark-3;
    width: 100%;
    font-size: 1.6rem;
    transition: all 0.3s;

    &:hover,
    &:focus {
      outline: none;
    }
  }

  &__button {
    text-align: center;
    background: $color-grey-light-1;
    border-radius: 5px;

    &:hover,
    &:active {
      background: $color-grey-light-2;
    }
  }
}
</style>
