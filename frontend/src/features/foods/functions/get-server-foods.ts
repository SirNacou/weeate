import { getFoods } from "@/client";
import { createServerFn } from "@tanstack/react-start";

export const getServerFoods = createServerFn({
  method: "GET",
}).handler(async () => {
  const res = await getFoods();
  return res.data;
});
