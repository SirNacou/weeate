import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  Card,
  CardContent,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { FoodImage } from "@/features/foods/components/food-image";

import { Button } from "@/components/animate-ui/components/buttons/button";
import {
  Dialog,
  DialogClose,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTrigger,
} from "@/components/animate-ui/components/radix/dialog";
import { ImageWithFallback } from "@/components/image-with-fallback";

const QRPayDialog = () => {
  return (
    <Dialog>
      <DialogTrigger asChild>
        <Button>Pay now</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <h3 className="text-lg font-semibold">Scan to Pay</h3>
        </DialogHeader>
        <div>
          <ImageWithFallback src={""} alt={"QR Code"}></ImageWithFallback>
        </div>
        <DialogFooter>
          <DialogClose>
            <Button variant="outline">Close</Button>
          </DialogClose>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
};

type Props = {};

const MealCard = ({}: Props) => {
  return (
    <div className="flex flex-col gap-3">
      <h3 className="text-2xl font-bold">Your Meal</h3>
      <Card>
        <CardHeader>
          <CardTitle className="flex justify-between text-xl">
            <span>Food Name</span>
            <span>
              {Intl.NumberFormat("vi-VN", {
                style: "currency",
                currency: "VND",
              }).format(100000)}
            </span>
          </CardTitle>
        </CardHeader>
        <CardContent className="p-0 aspect-square sm:aspect-video ">
          <FoodImage className="h-auto max-w-full" src={""} alt={""} />
        </CardContent>

        <CardFooter className="flex justify-between">
          <div className="flex items-center gap-2">
            <Avatar className="size-8 sm:size-10">
              <AvatarImage src={""} />
              <AvatarFallback>BY</AvatarFallback>
            </Avatar>
            <h3 className="text-base sm:text-lg leading-tight text-primary">
              Buyer
            </h3>
          </div>

          <QRPayDialog />
        </CardFooter>
      </Card>
    </div>
  );
};

export default MealCard;
