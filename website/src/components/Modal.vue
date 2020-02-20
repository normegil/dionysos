<template>
  <transition name="fade">
    <div v-if="show" class="modal__wrapper">
      <div class="modal__container">
        <div class="modal__header">
          <h1 class="modal__title">{{ title }}</h1>
          <Button icon="las la-times" @click="$emit('close')" />
        </div>
        <div class="modal__content">
          <slot />
        </div>
        <div class="modal__actions">
          <slot name="actions" />
        </div>
      </div>
    </div>
  </transition>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import { Prop } from "vue-property-decorator";
import Button from "./Button.vue";
@Component({
  components: { Button }
})
export default class Modal extends Vue {
  @Prop({ required: true })
  show!: boolean;

  @Prop({ required: true })
  title!: boolean;
}
</script>

<style lang="scss">
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter,
.fade-leave-to {
  opacity: 0;
}

.modal {
  &__wrapper {
    position: fixed;
    top: 0;
    left: 0;
    height: 100vh;
    width: 100vw;
    background-color: rgba($color-grey-dark-3, 0.8);
    z-index: 3000;
    transition: all 0.3s;
  }

  &__container {
    @include center-block;
    background-color: $color-white;
    min-width: 25%;
    box-shadow: 0 2rem 4rem rgba($color-black, 0.8);
    z-index: 4000;
    overflow: hidden;
    padding: 1.5rem;
  }

  &__header {
    display: flex;
    align-items: center;
    padding-bottom: 5px;
    border-bottom: 2px solid $color-grey-dark;
  }

  &__title {
    flex-grow: 2;
  }

  &__content {
    margin-top: 15px;
  }

  &__actions {
    margin-top: 15px;
    text-align: right;
  }
}
</style>
