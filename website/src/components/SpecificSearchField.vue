<template>
  <InputField
    v-model="searched"
    :button="$t('ui.button.search')"
    @keyup-enter="search"
    @button-click="search"
  />
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";
import InputField from "./InputField.vue";
@Component({
  components: { InputField }
})
export default class SpecificSearchField extends Vue {
  @Prop({ default: "", required: false })
  placeholder!: string;

  searched = "";

  search(): void {
    this.$emit("search", this.searched.trim());
  }
}
</script>

<style lang="scss"></style>
