<template>
  <a class="button" @click="click">
    <i v-if="icon !== '' && icon != null" :class="icon" />
    <span v-if="titleKey !== '' && titleKey != null">{{ $t(titleKey) }}</span>
  </a>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component
export default class Button extends Vue {
  @Prop({ default: "", required: false })
  icon!: string;

  @Prop({ default: "", required: false })
  titleKey!: string;

  click(): void {
    this.$emit("click");
  }
}
</script>

<style lang="scss">
.button {
  display: inline-block;
  padding: 0.5rem;
  border-radius: 5px;
  border: 1px solid $color-grey-light-2;
  text-decoration: none;
  color: $color-grey-dark-3;
  transition: all 0.3s;
  white-space: nowrap;

  &:active,
  &:hover {
    background-color: $color-grey-light-2;
  }
}
</style>
