<template>
  <div class="search-field">
    <!--suppress HtmlFormInputWithoutLabel -->
    <input
      class="search-field__input"
      type="text"
      :placeholder="$t('ui.components.search-field.placeholder')"
      v-model.trim="searched"
      @keyup.enter="search"
    />
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";

@Component
export default class SearchField extends Vue {
  searched = "";

  search(): void {
    console.log("Searching: " + this.searched);
  }
}
</script>

<style lang="scss">
.search-field {
  &__input {
    background-color: $color-white;
    display: inline-block;
    width: 100%;
    font-size: 1.5rem;
    color: $color-grey-dark-3;
    padding: 1.5rem;
    border: none;
    border-bottom: 1px solid $color-grey-light-2;
  }
}
</style>
