import { serverClient } from "@/api";
import { getFoods } from "@/client";
import { createServerFn } from "@tanstack/react-start";

export const getServerFoods = createServerFn({
  method: "GET",
}).handler(async () => {
  const res = await getFoods({ client: serverClient });
  console.log("Server Foods:", res);
  return res.data;
});
