import { listFoods } from "@/api"
import { zListFoodsData } from "@/client/zod.gen"
import { createServerFn } from "@tanstack/react-start"

export const getFoodsServer = createServerFn({
  method: "GET",
})
  .inputValidator(zListFoodsData)
  .handler(async ({ data }) => {
    const res = await listFoods({ ...data });
    console.log("getFoodsServer res:", res);
    if (res.error) {
      throw res.error;
    }
    return res.data;
  });
