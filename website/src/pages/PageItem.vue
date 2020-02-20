<template>
  <div class="content">
    <SearchField />
    <div class="content__main">
      <div class="content__title-container">
        <h1 class="content__title-text">{{ $t("ui.menu.main.items") }}</h1>
      </div>
      <div class="content__container flex-content">
        <Pagination
          :current-page="currentPage"
          :item-per-page="itemsPerPage"
          :number-of-pages="numberOfPages"
          @change-page="changePage"
        />
        <SpecificSearchField
          class="content__search-field"
          :placeholder="
            $t('ui.components.specific-search-field.items.placeholder')
          "
          :searched="searched"
          @search="filter"
        />
        <Button
          icon="las la-plus"
          titleKey="ui.button.add"
          @click="showUpdateItem = true"
        />
      </div>
      <div class="content__container">
        <Table :headings="headings" :content="items" @remove="removeItem" />
      </div>
    </div>
    <Modal :show="showUpdateItem" :title="modalLabel" @close="closeModal">
      <div class="item-modal__content">
        <InputField
          class="item-modal__field"
          v-model="editedModel.id"
          v-if="editedModel.id !== ''"
          :label="$t('ui.modal.item.id')"
          :disabled="true"
        />
        <InputField
          class="item-modal__field"
          v-model="editedModel.name"
          :label="$t('ui.modal.item.name')"
        />
      </div>
      <template v-slot:actions>
        <Button :title="$t('ui.button.save')" @click="saveItem" />
      </template>
    </Modal>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import SearchField from "../components/SearchField.vue";
import Table from "../components/Table.vue";
import Item from "../model/Item";
import Button from "../components/Button.vue";
import Pagination from "../components/Pagination.vue";
import SpecificSearchField from "../components/SpecificSearchField.vue";
import TableColumn from "../model/TableColumn";
import Modal from "../components/Modal.vue";
import InputField from "../components/InputField.vue";
@Component({
  components: {
    InputField,
    Modal,
    SpecificSearchField,
    Pagination,
    Button,
    Table,
    SearchField
  }
})
export default class PageItem extends Vue {
  showUpdateItem = false;

  editedModel: { id: string; name: string } = { id: "", name: "" };

  get items(): Item[] {
    return this.$store.state.items.items;
  }

  get searched(): string {
    return this.$store.state.items.filter;
  }

  get modalLabel(): string {
    let key = "ui.modal.item.title.add";
    if (this.editedModel.id !== "") {
      key = "ui.modal.item.title.update";
    }
    return this.$t(key) as string;
  }

  get numberOfPages(): number {
    return this.$store.getters["items/numberOfPages"];
  }

  get currentPage(): number {
    return this.$store.getters["items/currentPage"];
  }

  get itemsPerPage(): number {
    return this.$store.state.items.itemsPerPage;
  }

  get headings(): TableColumn[] {
    return [
      {
        name: this.$t("ui.components.table.items.heading.name") as string
      }
    ];
  }

  mounted(): void {
    this.$store.dispatch("items/load");
  }

  changePage(pageNb: number): void {
    this.$store.dispatch("items/changePage", pageNb);
  }

  filter(searchedValue: string): void {
    this.$store.dispatch("items/filter", searchedValue);
  }

  removeItem(id: string): void {
    this.$store.dispatch("items/delete", id);
  }

  saveItem(): void {
    this.$store.dispatch("items/save", this.editedModel).finally(() => {
      this.closeModal();
    });
  }

  closeModal(): void {
    this.editedModel = {
      id: "",
      name: ""
    };
    this.showUpdateItem = false;
  }
}
</script>

<style lang="scss">
.content {
  background-color: $color-white;

  overflow: auto;

  &__main {
    margin: 0 3.7rem;
  }

  &__title {
    &-text {
      font-family: "Rokkitt", "Ubuntu Mono", sans-serif;
      font-weight: 400;
      margin-left: 2rem;
      margin-top: 2.2rem;
      font-size: 2.8rem;
    }
  }

  &__container {
    margin-top: 2.2rem;
    padding: 1rem;
    border-radius: 2px;
    box-shadow: 0 0 5px $color-grey-light-2;
  }
  &__search-field {
    width: 100%;
    margin: 0 1rem;
  }
}
.flex-content {
  display: flex;
}
.item-modal {
  &__content {
    padding: 0 1rem;
  }

  &__field:not(:last-child) {
    margin-bottom: 5px;
  }
}
</style>
