import { serverClient } from "@/api";
import { listFoods } from "@/client";
import { zListFoodsData } from "@/client/zod.gen";
import { createServerFn } from "@tanstack/react-start";

export const getFoodsServer = createServerFn({
  method: "GET",
})
  .inputValidator(zListFoodsData)
  .handler(async ({ data }) => {
    const res = await listFoods({ client: serverClient, ...data });
    if (res.error) {
      throw res.error;
    }
    return res.data;
  });
