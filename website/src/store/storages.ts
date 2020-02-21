import Item from "../model/Item";
import CollectionStore from "./CollectionStore";

interface StorageDTO {
  id: string;
  name: string;
}

function convert(dto: StorageDTO): Item {
  return new Item(dto.id, dto.name);
}

export const STORAGES = new CollectionStore<Item, StorageDTO>(
  "storages",
  convert
);
