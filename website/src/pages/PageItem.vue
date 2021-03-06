<template>
  <div class="content">
    <SearchField />
    <div class="content__main">
      <div class="content__title-container">
        <h1 class="content__title-text">{{ $t("ui.menu.main.items") }}</h1>
      </div>
      <CollectionManager
        class="content__collection"
        store-namespace="items"
        :table-headers="headings"
        :filter-placeholder="
          $t('ui.components.specific-search-field.items.placeholder')
        "
        :items-per-page="$store.state.items.itemsPerPage"
        @create-item="createItem"
        @edit-item="editItem"
        @remove-item="removeItem"
      />
    </div>
    <Modal :show="showUpdateItemModal" :title="modalLabel" @close="closeModal">
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
import Button from "../components/Button.vue";
import Pagination from "../components/Pagination.vue";
import SpecificSearchField from "../components/SpecificSearchField.vue";
import TableColumn from "../model/TableColumn";
import Modal from "../components/Modal.vue";
import InputField from "../components/InputField.vue";
import CollectionManager from "../components/CollectionManager.vue";

@Component({
  components: {
    InputField,
    Modal,
    SpecificSearchField,
    Pagination,
    Button,
    Table,
    SearchField,
    CollectionManager
  }
})
export default class PageItem extends Vue {
  showUpdateItemModal = false;

  editedModel: { id: string; name: string } = { id: "", name: "" };

  get modalLabel(): string {
    let key = "ui.modal.item.title.add";
    if (this.editedModel.id !== "") {
      key = "ui.modal.item.title.update";
    }
    return this.$t(key) as string;
  }

  get headings(): TableColumn[] {
    return [
      {
        name: this.$t("ui.components.table.items.heading.name") as string
      }
    ];
  }

  createItem(): void {
    this.editedModel = {
      id: "",
      name: ""
    };
    this.showUpdateItemModal = true;
  }

  editItem(id: string): void {
    let found = false;
    for (const item of this.$store.state.items.items) {
      if (item.id === id) {
        this.editedModel = item;
        found = true;
        break;
      }
    }
    if (found) {
      this.showUpdateItemModal = true;
    } else {
      console.error("Could not find item with ID: " + id);
    }
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
    this.showUpdateItemModal = false;
  }

  mounted(): void {
    this.$store.dispatch("items/load");
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

  &__collection {
    margin-bottom: 2rem;
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
