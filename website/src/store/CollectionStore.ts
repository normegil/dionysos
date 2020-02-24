import { ActionContext, Module } from "vuex";
import { RootState } from "./model";
import { HTTP } from "./http";
import { AxiosError, AxiosResponse } from "axios";

interface Collection<T> {
  items: T[];
  itemsPerPage: number;
  currentIndex: number;
  total: number;
  filter: string;
}

interface CollectionResponse<DTO> {
  items: DTO[];
  totalSize: number;
  offset: number;
  limit: number;
  filter: string;
}

export default class CollectionStore<FinalType, DTO>
  implements Module<Collection<FinalType>, RootState> {
  namespaced = true;
  state = {
    items: [],
    currentIndex: 0,
    itemsPerPage: 20,
    total: 4,
    filter: ""
  };

  webServicePrefix: string;
  convert: (dto: DTO) => FinalType;

  constructor(webServicePrefix: string, convert: (dto: DTO) => FinalType) {
    this.webServicePrefix = webServicePrefix;
    this.convert = convert;
  }

  mutations = {
    setItems: (state: Collection<FinalType>, items: FinalType[]): void => {
      state.items = items;
    },
    setTotal: (state: Collection<FinalType>, nbItems: number): void => {
      state.total = nbItems;
    },
    setItemsPerPage: (
      state: Collection<FinalType>,
      itemsPerPage: number
    ): void => {
      state.itemsPerPage = itemsPerPage;
    },
    setCurrentIndex: (
      state: Collection<FinalType>,
      currentIndex: number
    ): void => {
      state.currentIndex = currentIndex;
    },
    setFilter: (state: Collection<FinalType>, filter: string): void => {
      state.filter = filter;
    }
  };

  actions = {
    load: (ctx: ActionContext<Collection<FinalType>, RootState>): void => {
      let url =
        "/" +
        this.webServicePrefix +
        "?limit=" +
        ctx.state.itemsPerPage +
        ";offset=" +
        ctx.state.currentIndex;
      if (ctx.state.filter !== "") {
        url += ";filter=" + ctx.state.filter;
      }
      HTTP.get(url)
        .then((r: AxiosResponse<CollectionResponse<DTO>>) => {
          return ctx.dispatch("refreshItems", r.data);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    },
    save: (
      ctx: ActionContext<Collection<FinalType>, RootState>,
      item: DTO
    ): void => {
      HTTP.put("/" + this.webServicePrefix, item)
        .then(() => {
          return ctx.dispatch("load");
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    },
    delete: (
      ctx: ActionContext<Collection<FinalType>, RootState>,
      id: string
    ): void => {
      HTTP.delete("/" + this.webServicePrefix + "/" + id)
        .then(() => {
          return ctx.dispatch("load");
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    },
    refreshItems: (
      ctx: ActionContext<Collection<FinalType>, RootState>,
      data: CollectionResponse<DTO>
    ): void => {
      const itemCollection = data.items.map(
        (dto): FinalType => this.convert(dto)
      );
      ctx.commit("setItems", itemCollection);
      ctx.commit("setTotal", data.totalSize);
      ctx.commit("setItemsPerPage", data.limit);
      ctx.commit("setFilter", data.filter);
      ctx.commit("setCurrentIndex", data.offset);
    }
  };
}
