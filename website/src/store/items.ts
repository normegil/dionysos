import Item from "../model/Item";
import CollectionStore from "./CollectionStore";

interface ItemDTO {
  id: string;
  name: string;
}

function convert(dto: ItemDTO): Item {
  return new Item(dto.id, dto.name);
}

export const ITEMS = new CollectionStore<Item, ItemDTO>(
  "items",
  "item",
  convert
);
