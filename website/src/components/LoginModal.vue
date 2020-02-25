<template>
  <Modal :title="$t('ui.modal.login.title')" :show="show" @close="close">
    <div class="login-modal__content" @keyup.enter="login">
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
      <Button :title="$t('ui.button.login')" @click="login" />
    </template>
  </Modal>
</template>

<script lang="ts">
import Component from "vue-class-component";
import Vue from "vue";
import Modal from "./Modal.vue";
import InputField from "./InputField.vue";
import Button from "./Button.vue";
import { Prop } from "vue-property-decorator";
@Component({
  components: { Modal, InputField, Button }
})
export default class LoginModal extends Vue {
  @Prop({ required: true })
  show!: boolean;

  username = "";
  password = "";

  login() {
    console.log("Login: " + this.username + ":" + this.password);
  }

  close() {
    this.$emit("close");
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
