import PageItem from "../pages/PageItem.vue";
import PageStorage from "../pages/PageStorage.vue"

export const ROUTES = [
  {
    path: "/storages",
    component: PageStorage
  },
  {
    path: "/items",
    component: PageItem
  }
];
