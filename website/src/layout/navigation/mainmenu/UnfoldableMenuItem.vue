<template>
  <div @click="switchDisplaySubItems" class="menu-foldable-item">
    <i :class="icon" class="menu-foldable-item__icon" />
    <span class="menu-foldable-item__text"> {{ title }}</span>
    <div class="menu-foldable-item__fold-icon-container">
      <i
        :style="caretRotation"
        class="las la-caret-right menu-foldable-item__fold-icon"
      />
    </div>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import LinkWithIcon from "../../../model/LinkWithIcon";

@Component
export default class UnfoldableMenuItem extends Vue {
  @Prop({ default: "", required: true })
  icon!: string;

  @Prop({ default: "", required: true })
  title!: string;

  @Prop({ default: [], required: true })
  subItems!: LinkWithIcon[];

  turn = 0;

  caretRotation(): { transform: string } {
    return { transform: "rotate(" + this.turn + "deg)" };
  }

  switchDisplaySubItems(): void {
    if (this.turn === 0) {
      this.turn = 90;
    } else {
      this.turn = 0;
    }
  }
}
</script>

<style lang="scss">
.menu-foldable-item {
  display: block;
  font-size: 2.5rem;
  padding: 1rem 0;
  border-bottom: 1px solid $color-grey-dark-3;
  padding-left: 3rem;
  color: $color-grey-dark;
  text-decoration: none;
  position: relative;
  transition: all 0.3s;

  &:hover {
    background-color: $color-grey-dark-3;
    color: $color-white;
  }

  &__text {
    display: inline-block;
    margin-left: 1rem;
  }

  &__fold-icon-container {
    position: absolute;
    right: 3rem;
    top: 1rem;
    color: rgba($color-black, 0);
    transition: color 0.3s;
  }

  &:hover &__fold-icon-container {
    color: rgba($color-white, 1);
  }
}
</style>
