<template>
  <Modal :title="$t('ui.modal.login.title')" :show="show" @close="close">
    <div class="login-modal__content" @keyup.enter="signIn">
      <p v-if="showError" class="login-modal__error-text">
        {{ $t("ui.modal.login.error") }}
      </p>
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
import { AxiosError } from "axios";
@Component({
  components: { Modal, InputField, Button }
})
export default class LoginModal extends Vue {
  username = "";
  password = "";

  showError = false;

  get show(): boolean {
    return this.$store.state.auth.showLoginModal;
  }

  signIn(): void {
    this.$store
      .dispatch("auth/signIn", {
        username: this.username,
        password: this.password
      })
      .then(() => {
        this.showError = false;
      })
      .catch((err: AxiosError): void => {
        if (undefined !== err.response) {
          if (err.response.status === 401) {
            this.showError = true;
          }
        }
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

  &__error-text {
    border-radius: 10px;
    border: 1px solid $color-red;
    background-color: $color-red-light;
    padding: 0.3rem 1rem;
    margin-bottom: 1.5rem;
  }

  &__field:not(:last-child) {
    margin-bottom: 5px;
  }
}
</style>
