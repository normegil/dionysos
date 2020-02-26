<template>
  <div
    class="field__container"
    :class="{ 'field__container--disabled': disabled }"
  >
    <label
      :for="this._uid"
      class="field__label"
      :class="{ 'field__label--disabled': disabled }"
      v-if="labelEnabled"
      >{{ label }}</label
    >
    <input
      :id="this._uid"
      class="field__input"
      :class="{
        'left-rounded': !labelEnabled && !password,
        'right-rounded': !buttonEnabled && !password
      }"
      :value="value"
      :placeholder="placeholder"
      @input="$emit('input', $event.target.value)"
      @keyup.enter="$emit('keyup-enter')"
      :disabled="disabled"
      :type="fieldType"
    />
    <input
      v-if="password"
      type="checkbox"
      v-model="showPassword"
      id="password-checkbox"
      class="field__password-checkbox"
    /><label
      v-if="password"
      for="password-checkbox"
      class="field__password-checkbox-label"
      :class="{
        'right-rounded': !buttonEnabled,
        'las la-eye': !showPassword,
        'las la-eye-slash': showPassword
      }"
    ></label>
    <a
      v-if="buttonEnabled"
      class="field__button"
      @click.stop="$emit('button-click')"
    >
      {{ button }}
    </a>
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

  @Prop({ default: "", required: false })
  placeholder!: string;

  @Prop({ default: "", required: false })
  button!: string;

  @Prop({ default: false, required: false })
  disabled!: boolean;

  @Prop({ default: false, required: false })
  password!: boolean;

  showPassword = false;

  get fieldType(): string {
    if (this.password && !this.showPassword) {
      return "password";
    }
    return "text";
  }

  get labelEnabled(): boolean {
    return this.label !== "";
  }

  get buttonEnabled(): boolean {
    return this.button !== "";
  }
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

    color: $color-grey-dark-3;
    transition: all 0.3s;
    padding: 0 1rem;

    &:hover,
    &:focus {
      outline: none;
    }

    &--disabled {
      cursor: not-allowed;
    }
  }

  &__button {
    text-align: center;
    padding: 0.5rem 2rem;
    border-left: 1px solid $color-grey-light-2;
    transition: all 0.3s;

    &:hover,
    &:active {
      background: $color-grey-light-2;
    }
  }

  &__password-checkbox {
    display: none;
  }

  &__password-checkbox-label {
    background-color: #fcfcfc;
    padding: 0.6rem 1rem;

    &:hover,
    &:active {
      color: $color-primary;
    }
  }
}
.left-rounded {
  border-bottom-left-radius: 5px;
  border-top-left-radius: 5px;
}

.right-rounded {
  border-bottom-right-radius: 5px;
  border-top-right-radius: 5px;
}
</style>
