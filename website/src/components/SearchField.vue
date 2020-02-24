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
  get searched(): string {
    return this.$store.state.search;
  }

  set searched(searched: string) {
    this.$store.commit("setSearch", searched);
  }

  search(): void {
    this.$store.dispatch("search", this.$router);
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
    padding: 1.4rem;
    border: none;
    border-bottom: 1px solid $color-grey-light-2;
  }
}
</style>
