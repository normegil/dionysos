import PageItem from "../pages/PageItem.vue";
import PageStorage from "../pages/PageStorage.vue";
import PageSearch from "../pages/PageSearch.vue";

export const ROUTES = [
  {
    path: "/storages",
    component: PageStorage
  },
  {
    path: "/items",
    component: PageItem
  },
  {
    path: "/search",
    component: PageSearch
  }
];
