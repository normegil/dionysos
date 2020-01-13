<template>
  <Header />
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import Header from "./layout/header/Header.vue";
@Component({
  components: { Header }
})
export default class App extends Vue {}
</script>

<style lang="scss"></style>
