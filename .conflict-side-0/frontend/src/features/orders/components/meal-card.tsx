import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { FoodImage } from "@/features/foods/components/food-image";

import { OrderedItem } from "@/client";

type Props = {
  orderedItems: OrderedItem[];
};

const MealCard = ({ orderedItems }: Props) => {
  return (
    <div className="flex flex-col gap-3">
      <h3 className="text-2xl font-bold">Your Meal</h3>
      {orderedItems.map((item) => {
        return (
          <Card>
            <CardHeader>
              <CardTitle className="flex justify-between text-xl">
                <span>{item.food_name}</span>
              </CardTitle>
            </CardHeader>
            <CardContent className="p-0 aspect-square sm:aspect-video ">
              <FoodImage
                className="h-auto max-w-full"
                src={item.food_url}
                alt={item.food_name}
              />
            </CardContent>

            <CardFooter className="flex justify-between items-center">
              <div className="flex items-center gap-2 w-full">
                <Avatar className="size-8 sm:size-10">
                  <AvatarImage src={item.buyer.avatar_url} />
                  <AvatarFallback>
                    {item.buyer.display_name?.slice(0, 2) ?? "BU"}
                  </AvatarFallback>
                </Avatar>
                <h3 className="text-base sm:text-lg leading-tight text-primary">
                  {item.buyer.display_name}
                </h3>
              </div>

              <div className="flex flex-col sm:flex-row items-center justify-end w-full text-lg">
                <span className="font-bold mr-2">Total Cost:</span>
                <span className="font-medium text-slate-700">
                  {Intl.NumberFormat("vi-VN", {
                    style: "currency",
                    currency: "VND",
                  }).format(item.total_price)}
                </span>
              </div>
              {/* <QRPayDialog /> */}
            </CardFooter>
          </Card>
        );
      })}
    </div>
  );
};

export default MealCard;
