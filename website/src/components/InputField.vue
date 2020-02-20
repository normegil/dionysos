<template>
  <div
    class="field__container"
    :class="{ 'field__container--disabled': disabled }"
  >
    <label
      :for="this._uid"
      class="field__label"
      :class="{ 'field__label--disabled': disabled }"
      v-if="label !== ''"
      >{{ label }}</label
    >
    <input
      :id="this._uid"
      class="field__input"
      :value="value"
      @input="$emit('input', $event.target.value)"
      :disabled="disabled"
    />
  </div>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";

@Component
export default class InputField extends Vue {
  @Prop({ default: "", required: false })
  label!: string;

  @Prop({ required: true })
  value!: string;

  @Prop({ default: false, required: false })
  disabled!: boolean;
}
</script>

<style lang="scss">
.field {
  &__container {
    display: flex;
    flex-direction: row;
    border: 1px solid $color-grey-light-2;
    border-radius: 5px;
    font-size: 1.6rem;
    transition: all 0.3s;

    &:hover:not(&--disabled),
    &:focus:not(&--disabled) {
      border: 1px solid $color-grey-light-1;
    }

    &--disabled {
      cursor: not-allowed;
    }
  }
  &__label {
    text-align: center;
    padding: 0.5rem 2rem;
    border-right: 1px solid $color-grey-light-2;
    transition: all 0.3s;
    border-bottom-left-radius: 5px;
    border-top-left-radius: 5px;

    &:hover:not(&--disabled),
    &:active:not(&--disabled) {
      background: $color-grey-light-3;
    }

    &--disabled {
      cursor: not-allowed;
    }
  }
  &__input {
    flex-grow: 2;
    border: none;
    border-bottom-right-radius: 5px;
    border-top-right-radius: 5px;

    color: $color-grey-dark-3;
    transition: all 0.3s;
    padding-left: 1rem;

    &:hover,
    &:focus {
      outline: none;
    }

    &--disabled {
      cursor: not-allowed;
    }
  }
}
</style>
