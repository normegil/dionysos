<template>
  <div class="menu-foldable-item">
    <div
      @click="switchDisplaySubItems"
      class="menu-foldable-item__main"
      :class="borderClass"
    >
      <i :class="icon" class="menu-foldable-item__main-icon" />
      <span class="menu-foldable-item__main-text"> {{ title }}</span>

      <div class="menu-foldable-item__main-fold-icon-container">
        <i
          :style="caretRotation"
          class="las la-caret-right menu-foldable-item__main-fold-icon"
        />
      </div>
    </div>
    <transition :name="slideFadeType">
      <div
        v-if="unfold"
        class="menu-foldable-item__sub-item-container"
        :style="subItemsContainerPositionStyle"
      >
        <a
          class="menu-foldable-item__sub-item"
          :class="borderClass"
          :href="item.href"
          v-for="item in subItems"
          :key="item.title"
        >
          <i :class="item.icon" />
          <span class="menu-foldable-item__sub-item-text">{{
            item.title
          }}</span>
        </a>
      </div>
    </transition>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import LinkWithIcon from "../model/LinkWithIcon";

@Component
export default class UnfoldableMenuItem extends Vue {
  @Prop({ default: "", required: true })
  icon!: string;

  @Prop({ default: "", required: true })
  title!: string;

  @Prop({ default: "down" })
  direction!: string;

  @Prop({ default: [], required: true })
  subItems!: LinkWithIcon[];

  @Prop({ default: false, required: false })
  startUnfold!: boolean;

  unfold = this.startUnfold;

  get borderClass(): string {
    if (this.direction === "down") {
      return "menu-foldable-item--border-bottom";
    } else {
      return "menu-foldable-item--border-top";
    }
  }

  get slideFadeType(): string {
    const transition = "slide-fade";
    if (this.direction === "down") {
      return transition + "-down";
    } else {
      return transition + "-up";
    }
  }

  get caretRotation(): { transform: string; transition: string } {
    let rotation = 0;
    if (this.unfold) {
      if (this.direction === "down") {
        rotation = 90;
      } else {
        rotation = -90;
      }
    }
    return {
      transform: "rotate(" + rotation + "deg)",
      transition: "all 0.3s"
    };
  }

  get subItemsContainerPositionStyle(): { bottom: string } {
    if (this.direction === "down") {
      return { bottom: "0rem" };
    }

    const mainItemHeight = 5.2; // Main item height

    return {
      bottom: mainItemHeight + "rem"
    };
  }

  switchDisplaySubItems(): void {
    this.unfold = !this.unfold;
  }
}
</script>

<style lang="scss">
.menu-foldable-item {
  &__main {
    display: block;
    font-size: 2.5rem;
    padding: 1rem 0;
    padding-left: 3rem;
    color: $color-grey-dark;
    text-decoration: none;
    position: relative;
    transition: all 0.3s;

    &:hover {
      background-color: $color-grey-dark-3;
      color: $color-white;
    }

    &-text {
      display: inline-block;
      margin-left: 1rem;
    }

    &-fold-icon-container {
      position: absolute;
      right: 3rem;
      top: 1rem;
      color: rgba($color-black, 0);
      transition: all 0.3s;
    }

    &:hover &-fold-icon-container {
      color: rgba($color-white, 1);
    }
  }
  &__sub-item {
    &-container {
      position: absolute;
      width: 100%;
    }

    display: block;
    font-size: 2rem;
    padding: 1rem 0;
    padding-left: 5rem;
    color: $color-grey-dark;
    text-decoration: none;

    &:hover {
      background-color: $color-grey-dark-3;
      color: $color-white;
    }

    &-text {
      display: inline-block;
      margin-left: 0.8rem;
    }
  }
  &--border-top {
    border-top: 1px solid $color-grey-dark-3;
  }
  &--border-bottom {
    border-bottom: 1px solid $color-grey-dark-3;
  }
}

.slide-fade {
  &-up-enter-active,
  &-up-leave-active,
  &-down-enter-active-down,
  &-down-leave-active-down {
    transition: all 0.3s ease;
  }

  &-up,
  &-down {
    &-enter,
    &-leave-to {
      opacity: 0;
    }
  }

  &-up {
    &-enter,
    &-leave-to {
      transform: translateY(1rem);
    }
  }

  &-down {
    &-enter,
    &-leave-to {
      transform: translateY(-1rem);
    }
  }
}
</style>
