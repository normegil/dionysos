<template>
  <Modal :title="$t('ui.modal.login.title')" :show="show" @close="close">
    <div class="login-modal__content" @keyup.enter="signIn">
      <InputField
        class="login-modal__field"
        v-model="username"
        :label="$t('ui.modal.login.username')"
      />
      <InputField
        class="login-modal__field"
        v-model="password"
        :label="$t('ui.modal.login.password')"
        :password="true"
      />
    </div>
    <template v-slot:actions>
      <Button :title="$t('ui.button.login')" @click="signIn" />
    </template>
  </Modal>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import Modal from "./Modal.vue";
import InputField from "./InputField.vue";
import Button from "./Button.vue";
@Component({
  components: { Modal, InputField, Button }
})
export default class LoginModal extends Vue {
  username = "";
  password = "";

  get show(): boolean {
    return this.$store.state.auth.showLoginModal;
  }

  signIn(): void {
    this.$store.dispatch("auth/signIn", {
      username: this.username,
      password: this.password
    });
  }

  close(): void {
    this.$store.commit("auth/setShowLoginModal", false);
  }
}
</script>

<style lang="scss">
.login-modal {
  &__content {
    padding: 0 1rem;
  }

  &__field:not(:last-child) {
    margin-bottom: 5px;
  }
}
</style>
