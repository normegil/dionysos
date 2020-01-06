<template>
    <div id="app">
        <h1>My Todo App!</h1>
        <p>{{test}}</p>
    </div>
</template>

<script lang="ts">
    import Vue from 'vue'
    import Component from 'vue-class-component'

    @Component
    export default class App extends Vue {
        test: string = "TSTest";
    }
</script>

<style>
</style>